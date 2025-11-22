package handler

import (
	"errors"

	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler はHTTPリクエストを処理するのだ
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler は UserHandler のインスタンスを生成するのだ
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register はユーザー登録APIのエンドポイントなのだ
func (h *UserHandler) Register(c *gin.Context) {
	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Service を呼び出してビジネスロジックを実行
	createdUser, err := h.userService.RegisterUser(&req)
	if err != nil {
		// サービスから返されたエラーが「メール重複エラー」かチェック
		if errors.Is(err, service.ErrEmailConflict) {
			// HTTP 409 (Conflict) を返す
			c.JSON(http.StatusConflict, gin.H{"error:This email address is already in user": err.Error()})
			return
		}

		// その他のサーバー内部エラー
		c.JSON(http.StatusInternalServerError, gin.H{"error": "内部サーバーエラーが発生しました"})
		return
	}

	// 3. 成功レスポンス
	res := model.UserRegisterResponse{
		UserID:      createdUser.UserID,
		UserName:    createdUser.UserName,
		DisplayName: createdUser.DisplayName,
		MailAddress: createdUser.MailAddress,
		CreatedAt:   createdUser.CreatedAt,
	}
	c.JSON(http.StatusCreated, res)
}

// Login はログインAPIのエンドポイントなのだ
func (h *UserHandler) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.LoginUser(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	res := model.UserLoginResponse{
		Token: token,
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	// AuthMiddlewareでセットされた user_id を取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// userIDはTokenクレームから文字列として来ているはず
	user, err := h.userService.GetUser(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, user)
}
