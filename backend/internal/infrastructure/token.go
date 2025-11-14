package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken はユーザーIDを元にJWT（通行証）を生成するのだ
func GenerateToken(userID string) (string, error) {
	// JWTに含める情報（Claims）
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 有効期限: 24時間

	// ヘッダーとペイロードを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵（.envから取得）で署名する
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET が設定されていません")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("トークンの署名に失敗しました: %v", err)
	}

	return tokenString, nil
}
