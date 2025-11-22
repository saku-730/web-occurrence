package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// CouchDBClient はCouchDBとの内部通信を行うのだ
type CouchDBClient interface {
	GetSessionCookie(username string) (string, error)
	CreateCouchDBUser(username string, password string) error
	UpsertDocument(docID string, data map[string]interface{}) error // 追加したのだ
}

type couchDBClient struct {
	client    *http.Client
	baseURL   string
	dbName    string // DB名も保持しておくのだ
	adminUser string
	adminPass string
}

// NewCouchDBClient は CouchDBClient のインスタンスを生成するのだ
func NewCouchDBClient() CouchDBClient {
	dbName := os.Getenv("COUCHDB_DB_NAME")
	if dbName == "" {
		dbName = "occurrence" // 環境変数がなければデフォルト値なのだ
	}

	return &couchDBClient{
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   os.Getenv("COUCHDB_URL"),
		dbName:    dbName,
		adminUser: os.Getenv("COUCHDB_ADMIN_USER"),
		adminPass: os.Getenv("COUCHDB_ADMIN_PASS"),
	}
}

// GetSessionCookie は、セッション代理発行を実行するのだ
func (c *couchDBClient) GetSessionCookie(username string) (string, error) {
	reqBody := map[string]string{"name": username}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("CouchDBリクエストボディのJSON化に失敗: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/_session", c.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("CouchDBリクエストの作成に失敗: %v", err)
	}

	req.SetBasicAuth(c.adminUser, c.adminPass)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("CouchDBへのセッション要求に失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CouchDBがセッション発行に失敗 (ステータス: %d)", resp.StatusCode)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "AuthSession" {
			return cookie.String(), nil
		}
	}

	return "", fmt.Errorf("CouchDBのレスポンスに AuthSession が見つかりません")
}

// CreateCouchDBUser は CouchDB の _users データベースにユーザーを新規作成する
func (c *couchDBClient) CreateCouchDBUser(username string, password string) error {
	reqBody := map[string]interface{}{
		"type":     "user",
		"name":     username,
		"password": password,
		"roles":    []string{},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("CouchDBユーザー作成リクエストのJSON化に失敗: %v", err)
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/_users/org.couchdb.user:%s", c.baseURL, username),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("CouchDBユーザー作成リクエストの作成に失敗: %v", err)
	}

	req.SetBasicAuth(c.adminUser, c.adminPass)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("CouchDBへのユーザー作成要求に失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("CouchDBがユーザー作成に失敗 (ステータス: %d)", resp.StatusCode)
	}

	return nil
}

// UpsertDocument はドキュメントを作成または更新するのだ
// 既存のドキュメントがある場合は _rev を取得して上書きするのだ
func (c *couchDBClient) UpsertDocument(docID string, data map[string]interface{}) error {
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, docID)

	// 1. まずGETして、既存の _rev があるか確認するのだ
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("GETリクエスト作成失敗: %v", err)
	}
	req.SetBasicAuth(c.adminUser, c.adminPass)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ドキュメント取得失敗: %v", err)
	}
	defer resp.Body.Close()

	// 200 OK なら既存データがあるので _rev を取得して data にマージするのだ
	if resp.StatusCode == http.StatusOK {
		var currentDoc map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&currentDoc); err != nil {
			return fmt.Errorf("既存ドキュメントのデコード失敗: %v", err)
		}
		if rev, ok := currentDoc["_rev"].(string); ok {
			data["_rev"] = rev
		}
	}

	// 2. PUT でデータを書き込むのだ
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON化失敗: %v", err)
	}

	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("PUTリクエスト作成失敗: %v", err)
	}
	req.SetBasicAuth(c.adminUser, c.adminPass)
	req.Header.Set("Content-Type", "application/json")

	resp, err = c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ドキュメント保存リクエスト失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("ドキュメント保存失敗 (ステータス: %d)", resp.StatusCode)
	}

	return nil
}
