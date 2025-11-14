package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword はbcryptを使ってパスワードをハッシュ化するのだ
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14はコスト
	return string(bytes), err
}

// CheckPasswordHash はハッシュと平文のパスワードを比較するのだ
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // エラーがなければ (nil なら) true (一致)
}
