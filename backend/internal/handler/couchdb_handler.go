package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/web-occurrence/backend/internal/service"
)

type CouchDBHandler struct {
	couchDBService service.CouchDBService
}

func NewCouchDBHandler(couchDBService service.CouchDBService) *CouchDBHandler {
	return &CouchDBHandler{couchDBService: couchDBService}
}

// GetCouchDBSession はCookie認証用のセッションを発行する場合に使用（今回はProxyメインなら使わない可能性あり）
func (h *CouchDBHandler) GetCouchDBSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id missing"})
		return
	}
	cookieString, err := h.couchDBService.RequestCouchDBSession(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Set-Cookie", cookieString)
	c.JSON(http.StatusOK, gin.H{"message": "Session Created"})
}

// ProxyRequest はフロントエンドからのPouchDB同期リクエストをCouchDBへ中継する
func (h *CouchDBHandler) ProxyRequest(c *gin.Context) {
	// Debugログ
	fmt.Println("--- [DEBUG] 1. ProxyRequest Handler Reached ---")

	// 1. ミドルウェアで認証済みの UserID を取得
	// AuthMiddlewareがセットした "user_id" (string) を取得
	userIDVal, exists := c.Get("user_id")
	if !exists {
		fmt.Println("--- [DEBUG] Error: user_id not found in context ---")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}
	userID := userIDVal.(string)
	fmt.Printf("--- [DEBUG] 2. UserID found: %s ---\n", userID)

	// 2. Serviceを使って、CouchDB用のユーザー名と署名トークンを取得
	// CouchDBのProxy Authentication設定で有効なHMACトークンを生成
	username, token, err := h.couchDBService.GenerateProxyCredentials(userID)
	if err != nil {
		fmt.Printf("--- [DEBUG] Error: GenerateProxyCredentials failed: %v ---\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報の生成に失敗しました: " + err.Error()})
		return
	}
	fmt.Printf("--- [DEBUG] 3. Credentials Generated. User: %s ---\n", username)

	// 3. 転送先URLの準備 (CouchDBのベースURL)
	targetURLStr := h.couchDBService.GetCouchDBURL()
	target, err := url.Parse(targetURLStr)
	if err != nil {
		fmt.Printf("--- [DEBUG] Error: URL Parse failed: %v ---\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CouchDBのURL設定が不正です"})
		return
	}

	// 4. リバースプロキシの作成
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Director: リクエスト内容を書き換える関数
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// パスの調整: /api/couchdb プレフィックスを削除して、CouchDBへのパスにする
		// 例: /api/couchdb/test_db/doc1 -> /test_db/doc1
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/couchdb")
		
		// 認証ヘッダーの注入 (Proxy Authentication)
		req.Header.Set("X-Auth-CouchDB-UserName", username)
		req.Header.Set("X-Auth-CouchDB-Roles", "member") // 必要に応じてロールを調整
		req.Header.Set("X-Auth-CouchDB-Token", token)

		// ホストヘッダーの書き換え (バックエンド側のHostに合わせる)
		req.Host = target.Host

		// デバッグ出力
		fmt.Println("--- [DEBUG] 5. Proxy Director executed ---")
		fmt.Printf("Target Path: %s\n", req.URL.Path)
		fmt.Printf("Header [UserName]: %s\n", req.Header.Get("X-Auth-CouchDB-UserName"))
		// トークンはセキュリティのためログには出さないか、一部伏せ字推奨
		// fmt.Printf("Header [Token]: %s\n", req.Header.Get("X-Auth-CouchDB-Token"))
		fmt.Println("-------------------------------------")
	}

	// エラーハンドリング
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("--- [DEBUG] Proxy Error: %v ---\n", err)
		http.Error(w, "Proxy Error: "+err.Error(), http.StatusBadGateway)
	}

	// 5. プロキシ実行
	fmt.Println("--- [DEBUG] 6. Starting ServeHTTP ---")
	proxy.ServeHTTP(c.Writer, c.Request)
	
	// Ginの処理をここで終了（Proxyがレスポンスを書き込むため）
	c.Abort()
}
