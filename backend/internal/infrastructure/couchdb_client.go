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
	CreateCouchDBUser(username string, password string) error // (今回追加)
}

type couchDBClient struct {
	client    *http.Client
	baseURL   string
	adminUser string
	adminPass string
}

// NewCouchDBClient は CouchDBClient のインスタンスを生成するのだ
func NewCouchDBClient() CouchDBClient {
	return &couchDBClient{
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   os.Getenv("COUCHDB_URL"),
		adminUser: os.Getenv("COUCHDB_ADMIN_USER"),
		adminPass: os.Getenv("COUCHDB_ADMIN_PASS"),
	}
}

// GetSessionCookie は、セッション代理発行を実行するのだ (既存)
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
	
	// 1. CouchDBが要求するユーザードキュメントの形式でJSONを作成
	// (パスワードは平文で送ると、CouchDBが自動でハッシュ化してくれるのだ)
	reqBody := map[string]interface{}{
		"type":     "user",
		"name":     username,
		"password": password,
		"roles":    []string{}, // ここでは空。必要なら "user_role" など追加
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("CouchDBユーザー作成リクエストのJSON化に失敗: %v", err)
	}

	// 2. _users DBのエンドポイント (org.couchdb.user:[username]) に `PUT` リクエストを作成
	// (CouchDBのユーザー作成は `POST /_users` ではなく `PUT /_users/[doc_id]` が推奨なのだ)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/_users/org.couchdb.user:%s", c.baseURL, username),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("CouchDBユーザー作成リクエストの作成に失敗: %v", err)
	}

	// 3. 「管理者」として認証
	req.SetBasicAuth(c.adminUser, c.adminPass)
	req.Header.Set("Content-Type", "application/json")

	// 4. リクエスト送信
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("CouchDBへのユーザー作成要求に失敗: %v", err)
	}
	defer resp.Body.Close()

	// 5. 201 (Created) または 202 (Accepted) 以外はエラーとする
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		// (409 Conflict なら「既にユーザーが存在する」エラー)
		return fmt.Errorf("CouchDBがユーザー作成に失敗 (ステータス: %d)", resp.StatusCode)
	}

	return nil
}
