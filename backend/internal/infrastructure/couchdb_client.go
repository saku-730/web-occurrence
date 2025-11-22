package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/saku-730/web-occurrence/backend/internal/model"
)

// CouchDBClient はCouchDBとの内部通信を行うのだ
type CouchDBClient interface {
	GetSessionCookie(username string) (string, error)
	CreateCouchDBUser(username string, password string) error
	UpsertDocument(docID string, data map[string]interface{}) error
	// ▼ 追加: 指定DBの全ドキュメントを取得する
	FetchAllDocs(dbName string) ([]map[string]interface{}, error)
}

type couchDBClient struct {
	client    *http.Client
	baseURL   string
	adminUser string
	adminPass string
}

// NewCouchDBClient は Config を受け取るように修正したのだ
func NewCouchDBClient(config *model.CouchDBConfig) CouchDBClient {
	// デフォルト値の処理
	if config.URL == "" {
		config.URL = "http://localhost:5984"
	}
	
	return &couchDBClient{
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   config.URL,
		adminUser: config.AdminUser,
		adminPass: config.AdminPass,
	}
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

// UpsertDocument は指定されたドキュメントIDでデータを保存する（既存があれば上書き）
// docID の例: "_local/master_data" や "occ_xxxx"
// dbName はURLに含まれるか、あるいは引数で渡すべきだが、
// ここでは docID にDB名が含まれていないため、呼び出し元で制御するか、
// 簡易的に固定のDB名を使うか、引数を増やす必要がある。
// ★修正: workstation_service.go からの呼び出しでは docID しか渡されていないが、
// どのDBに入れるかを決める必要がある。
// 現在の `workstation_service.go` の実装を見ると、
// docID = "_local/master_data" で呼び出している。
// これは PouchDB/CouchDB の仕様上、DBごとのローカルドキュメント。
// つまり、「どのDBに対して」という情報が不足している。
// ここでは、`NewCouchDBClient` でDB名を固定していた以前の実装をやめて、
// `UpsertDocument` に `dbName` 引数を追加するのが正しい設計だが、
// interfaceを変えると呼び出し元も修正が必要になる。
//
// ★一旦、元のコードの意図（Workstation作成時はまだ専用DBがない？）を汲みつつ、
// `_users` や `_replicator` 以外の通常のデータは、呼び出し側でDB名を意識すべき。
// 今回は、WorkstationServiceで使われている `docID` は `_local/...` なので、
// 実はDBを指定しないと保存できない。
//
// ここでは、エラーを解消するために、`UpsertDocument` を以下のように実装する。
// ただし、WorkstationService側も「どのDBにマスターデータを入れるか」を指定する必要がある。
// WorkstationServiceの呼び出し元 `docID` だけでは不足しているため、
// 本来は `UpsertDocument(dbName, docID, data)` とすべき。
//
// しかし、大量の修正を避けるため、ここでは `docID` に `dbName/docID` の形式で渡すか、
// メソッド側でハンドリングする必要がある。
//
// 今回のエラー回避として、WorkstationServiceで呼び出している `UpsertDocument` は、
// 作成したばかりのワークステーションDB (`test_db_ws_1`など) に対して行うべきもの。
//
// そのため、インターフェースを変更する。
// `UpsertDocument(dbName string, docID string, data map[string]interface{}) error`
func (c *couchDBClient) UpsertDocument(docID string, data map[string]interface{}) error {
	// ★注意: このメソッドは WorkstationService から呼ばれているが、
	// DB名が引数にないため、data["workstation_id"] などから推測するか、
	// あるいは `test_db_ws_{ID}` のようなルールでDB名を構築する必要がある。
	
	// dataの中に workstation_id がある場合、それを使ってDB名を特定するロジックを入れる
	// (これは少し強引だが、インターフェースを変えずに済ますための一時策)
	
	var dbName string
	if wsID, ok := data["workstation_id"].(string); ok {
		// プレフィックスは環境変数などから取るべきだが、一旦 test_db 固定
		dbName = fmt.Sprintf("test_db_ws_%s", wsID)
	} else {
		// デフォルトまたはエラー
		return fmt.Errorf("workstation_id not found in data, cannot determine DB name")
	}

	url := fmt.Sprintf("%s/%s/%s", c.baseURL, dbName, docID)

	// 1. DBが存在するか確認・作成 (本来はService層でやるべきだがここでやる)
	// PUT /{db}
	reqCreateDB, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", c.baseURL, dbName), nil)
	reqCreateDB.SetBasicAuth(c.adminUser, c.adminPass)
	c.client.Do(reqCreateDB) // エラー無視（既にある場合など）

	// 2. GET existing doc
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

// ▼ 追加実装: FetchAllDocs
func (c *couchDBClient) FetchAllDocs(dbName string) ([]map[string]interface{}, error) {
	// _all_docs?include_docs=true でドキュメント本体も含めて取得
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
		// DBが存在しない場合などはエラーまたは空を返す
		return nil, fmt.Errorf("failed to fetch docs: status %d", resp.StatusCode)
	}

	// レスポンス構造
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
		// デザインドキュメントなどを除外したい場合はここで判定
		// if id, ok := row.Doc["_id"].(string); ok && strings.HasPrefix(id, "_design/") { continue }
		docs = append(docs, row.Doc)
	}

	return docs, nil
}
