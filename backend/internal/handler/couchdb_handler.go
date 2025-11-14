package handler

import (
	"github.com/saku-730/web-occurrence/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CouchDBHandler struct {
	couchDBService service.CouchDBService
}

// NewCouchDBHandler は CouchDBHandler のインスタンスを生成するのだ
func NewCouchDBHandler(couchDBService service.CouchDBService) *CouchDBHandler {
	return &CouchDBHandler{couchDBService: couchDBService}
}

// GetCouchDBSession は GET /api/couchdb-session のエンドポイントなのだ
func (h *CouchDBHandler) GetCouchDBSession(c *gin.Context) {

	// 1. ミドルウェアが検証・セットした user_id を取得
	// (c.GetString は、キーが存在しないと空文字とfalseを返す)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ミドルウェアから user_id が渡されませんでした"})
		return
	}

	// 2. Service を呼び出してセッションCookie文字列を取得
	cookieString, err := h.couchDBService.RequestCouchDBSession(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. sのフローの「ステップ4」:
	// GoサーバーがCouchDBから受け取ったCookie文字列を、
	// Next.js（ブラウザ）へのレスポンスヘッダーに "Set-Cookie" としてセットする
	c.Header("Set-Cookie", cookieString)
	
	// (セキュリティのための追加設定)
	// Cookieを httpOnly (JSから読めなくする) にしたり、
	// Secure (HTTPSのみ) にしたりする設定もここで行うのが望ましいのだ
	// c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, true, true)

	// 4. 成功レスポンスを返す
	c.JSON(http.StatusOK, gin.H{"message": "CouchDBセッションが発行されました"})
}
