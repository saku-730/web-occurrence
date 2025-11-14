package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword はbcryptを使ってパスワードをハッシュ化するのだ
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14はコスト
	return string(bytes), err
}
