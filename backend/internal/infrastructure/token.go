package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() 

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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

func ValidateAndExtractUserID(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET が設定されていません")
	}

	// トークンをパース（解析）する
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名方法がHS256であることを確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("予期しない署名方法です: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err // (有効期限切れなどもここに含まれる)
	}

	// トークンが有効で、中身（Claims）がMapClaims形式かチェック
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// "user_id" を取り出す
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
	}

	return "", fmt.Errorf("無効なトークンです")
}
