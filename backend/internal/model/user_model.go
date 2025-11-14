package model

import "time"

// UserRegisterRequest はユーザー登録APIのJSONリクエストボディなのだ
type UserRegisterRequest struct {
	MailAddress string `json:"mailaddress" binding:"required,email"` // バリデーション
	Password    string `json:"password" binding:"required,min=8"`    // 8文字以上
}

// UserRegisterResponse はユーザー登録APIの成功レスポンスなのだ
type UserRegisterResponse struct {
	UserID      int64    `json:"user_id"`
	UserName    string    `json:"user_name"`
	DisplayName string    `json:"display_name"`
	MailAddress string    `json:"mail_address"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserLoginRequest はログインAPIのJSONリクエストボディなのだ
type UserLoginRequest struct {
	MailAddress string `json:"mailaddress" binding:"required,email"`
	Password    string `json:"password" binding:"required"` // ログイン時はminチェック不要
}

// UserLoginResponse はログインAPIの成功レスポンスなのだ
type UserLoginResponse struct {
	Token string `json:"token"` // JWTトークンを返す
}
