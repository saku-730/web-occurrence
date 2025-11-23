
package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strconv"

	"github.com/saku-730/web-occurrence/backend/internal/model"
)

type CouchDBClient interface {
	GetSessionCookie(username string) (string, error)
	CreateCouchDBUser(username string, password string) error
	UpsertDocument(docID string, data map[string]interface{}) error
	FetchAllDocs(dbName string) ([]map[string]interface{}, error)
	CreateDatabase(dbName string) error
	// ▼ 追加: ワークステーションIDからDB名を生成するヘルパーなのだ
	CreateWorkstationDBName(workstationID int64) string
	// ▼ 追加: DBにアクセス権を設定するメソッドなのだ
	SetDatabaseUserAccess(dbName string, userID string) error
}

type couchDBClient struct {
	client    *http.Client
	baseURL   string
	adminUser string
	adminPass string
	dbPrefix  string // DBプレフィックスを追加したのだ
}

func NewCouchDBClient(config *model.CouchDBConfig) CouchDBClient {
	if config.URL == "" {
		config.URL = "http://localhost:5984"
	}
	
	// ★修正: DBプレフィックスを "db" に修正するのだ
	dbPrefix := "db" // 環境変数からの取得が未実装のため一旦固定

	return &couchDBClient{
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   config.URL,
		adminUser: config.AdminUser,
		adminPass: config.AdminPass,
		dbPrefix:  dbPrefix,
	}
}

// ▼ 追加実装: ワークステーションIDからDB名を生成するのだ
func (c *couchDBClient) CreateWorkstationDBName(workstationID int64) string {
	// 修正後の命名規則: db_ws_[ID] になるのだ
	return fmt.Sprintf("%s_ws_%d", c.dbPrefix, workstationID)
}

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

// CreateDatabase は、管理者権限を使って新しいDBを作成するのだ
func (c *couchDBClient) CreateDatabase(dbName string) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, dbName)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("DB作成リクエスト作成失敗: %w", err)
	}
	// ★ここが認証の鍵！管理者情報を使うのだ
	req.SetBasicAuth(c.adminUser, c.adminPass)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("CouchDBへの接続失敗: %w", err)
	}
	defer resp.Body.Close()
	
	// 412 Precondition Failed (既に存在する) は成功とみなす
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusPreconditionFailed {
		return nil
	}

	// 401 が出たら、管理者認証が失敗しているのだ
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("DB作成失敗: 401 Unauthorized. 管理者ユーザー名/パスワードを確認してください")
	}

	return fmt.Errorf("DB作成失敗 (ステータス: %d)", resp.StatusCode)
}

// ▼ 追加実装: SetDatabaseUserAccess は指定したユーザーにDBの読み書き権限を与えるのだ
func (c *couchDBClient) SetDatabaseUserAccess(dbName string, userID string) error {
	securityDoc := map[string]interface{}{
		// ★ここが重要: ユーザーIDをDBのメンバーに追加することで403を解消するのだ
		"members": map[string]interface{}{
			"names": []string{userID}, // ユーザーIDをメンバーに追加するのだ
			"roles": []string{},
		},
		"admins": map[string]interface{}{
			"names": []string{},
			"roles": []string{"_admin"}, // 管理者ロールはそのまま残すのだ
		},
	}
	
	jsonData, err := json.Marshal(securityDoc)
	if err != nil {
		return fmt.Errorf("セキュリティドキュメントのJSON化に失敗: %w", err)
	}

	url := fmt.Sprintf("%s/%s/_security", c.baseURL, dbName)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("セキュリティ設定リクエスト作成失敗: %w", err)
	}
	
	// 管理者認証を使ってセキュリティ設定を更新するのだ
	req.SetBasicAuth(c.adminUser, c.adminPass)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("CouchDBへのセキュリティ設定要求に失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("セキュリティ設定失敗 (ステータス: %d)", resp.StatusCode)
	}

	return nil
}

// UpsertDocument はドキュメントを作成または更新するのだ
func (c *couchDBClient) UpsertDocument(docID string, data map[string]interface{}) error {
	// data から workstation_id (string) を取得してDB名を決定するのだ
	wsIDStr, ok := data["workstation_id"].(string)
	if !ok {
		return fmt.Errorf("workstation_id がデータに見つからないか、string型ではありません")
	}

	// WorkstationIDをint64に変換
	wsID, err := strconv.ParseInt(wsIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("workstation_id の数値変換に失敗: %w", err)
	}

	dbName := c.CreateWorkstationDBName(wsID)
	
	// ★修正: DB作成を試行し、エラーをチェックするのだ！
	// ここで401エラーが出れば、次のドキュメント保存を試行せずに即座にエラーを返すのだ
	if err := c.CreateDatabase(dbName); err != nil {
		return fmt.Errorf("DBの確保に失敗 (%s): %w", dbName, err)
	}

	url := fmt.Sprintf("%s/%s/%s", c.baseURL, dbName, docID)

	// 2. GET existing doc (for _rev)
	// ... (ここはデータ取得ロジックなので省略、データ更新ロジックはそのまま)
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

	if resp.StatusCode == http.StatusOK {
		var currentDoc map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&currentDoc); err != nil {
			return fmt.Errorf("既存ドキュメントのデコード失敗: %v", err)
		}
		if rev, ok := currentDoc["_rev"].(string); ok {
			data["_rev"] = rev
		}
	}

	// 3. PUT doc
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
		return fmt.Errorf("ドキュメント保存失敗 (ステータス: %d) %s", resp.StatusCode, url)
	}

	return nil
}

func (c *couchDBClient) FetchAllDocs(dbName string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s/_all_docs?include_docs=true", c.baseURL, dbName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.adminUser, c.adminPass)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch docs: status %d", resp.StatusCode)
	}

	type Row struct {
		Doc map[string]interface{} `json:"doc"`
	}
	type AllDocsResponse struct {
		Rows []Row `json:"rows"`
	}

	var result AllDocsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var docs []map[string]interface{}
	for _, row := range result.Rows {
		docs = append(docs, row.Doc)
	}

	return docs, nil
}
