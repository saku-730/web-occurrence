package handler

import (
	"github.com/saku-730/web-occurrence/backend/internal/service"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CouchDBHandler struct {
	couchDBService service.CouchDBService
}

func NewCouchDBHandler(couchDBService service.CouchDBService) *CouchDBHandler {
	return &CouchDBHandler{couchDBService: couchDBService}
}

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

func (h *CouchDBHandler) ProxyRequest(c *gin.Context) {
	// 【チェックポイント1】 ハンドラーに入ったか確認
	fmt.Println("--- [DEBUG] 1. ProxyRequest Handler Reached ---")

	// 1. ミドルウェアで認証済みの UserID を取得
	userIDVal, exists := c.Get("user_id")
	if !exists {
		fmt.Println("--- [DEBUG] Error: user_id not found in context ---")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}
	userID := userIDVal.(string)
	fmt.Printf("--- [DEBUG] 2. UserID found: %s ---\n", userID)

	// 2. Serviceを使って、CouchDB用のユーザー名と署名トークンを取得
	username, token, err := h.couchDBService.GenerateProxyCredentials(userID)
	if err != nil {
		fmt.Printf("--- [DEBUG] Error: GenerateProxyCredentials failed: %v ---\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報の生成に失敗しました: " + err.Error()})
		return
	}
	fmt.Printf("--- [DEBUG] 3. Credentials Generated. User: %s ---\n", username)

	// 3. 転送先URLの準備
	targetURLStr := h.couchDBService.GetCouchDBURL()
	fmt.Printf("--- [DEBUG] 4. Target URL string: %s ---\n", targetURLStr)
	
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

		// パスの調整
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/couchdb")
		
		// 認証ヘッダーの注入
		req.Header.Set("X-Auth-CouchDB-UserName", username)
		req.Header.Set("X-Auth-CouchDB-Roles", "member")
		req.Header.Set("X-Auth-CouchDB-Token", token)

		// ホストヘッダーの書き換え
		req.Host = target.Host

		// 【チェックポイント5】 ここが出なければ、プロキシ実行直前で何かが起きている
		fmt.Println("--- [DEBUG] 5. Proxy Director executed ---")
		fmt.Printf("Target Path: %s\n", req.URL.Path)
		fmt.Printf("Header [UserName]: %s\n", req.Header.Get("X-Auth-CouchDB-UserName"))
		fmt.Printf("Header [Token]: %s\n", req.Header.Get("X-Auth-CouchDB-Token"))
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
	
	c.Abort()
}
