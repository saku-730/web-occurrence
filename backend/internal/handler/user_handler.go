package handler

import (
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

	// 1. JSONリクエストを model.UserRegisterRequest にバインド（変換）＆バリデーション
	if err := c.ShouldBindJSON(&req); err != nil {
		// バリデーションエラー（email形式じゃない、パスワードが短いなど）
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Service を呼び出してビジネスロジックを実行
	createdUser, err := h.userService.RegisterUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 成功レスポンス（model.UserRegisterResponse）を作成
	res := model.UserRegisterResponse{
		UserID:      createdUser.UserID,
		UserName:    createdUser.UserName,
		DisplayName: createdUser.DisplayName,
		MailAddress: createdUser.MailAddress,
		CreatedAt:   createdUser.CreatedAt,
	}

	// 201 Created ステータスでレスポンスを返す
	c.JSON(http.StatusCreated, res)
}
