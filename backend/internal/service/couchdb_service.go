package service

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
)

// CouchDBService はCouchDB関連のビジネスロジックを担当するのだ
type CouchDBService interface {
	RequestCouchDBSession(userID string) (string, error)
	GenerateProxyCredentials(userID string) (string, string, error)
	GetCouchDBURL() string
}

type couchDBService struct {
	userRepo     repository.UserRepository
	couchClient  infrastructure.CouchDBClient
	configSecret string
	configURL    string
}

func NewCouchDBService(
	userRepo repository.UserRepository,
	couchClient infrastructure.CouchDBClient,
	secret string,
	url string,
) CouchDBService {
	return &couchDBService{
		userRepo:     userRepo,
		couchClient:  couchClient,
		configSecret: secret,
		configURL:    url,
	}
}

// RequestCouchDBSession (既存機能: 必要な場合のみ使用)
func (s *couchDBService) RequestCouchDBSession(userID string) (string, error) {
	// ★ここも user_id をそのまま使うように修正できるけど、
	// Cookie認証を使わないならこのメソッド自体削除してもいいのだ
	return s.couchClient.GetSessionCookie(userID)
}

// GenerateProxyCredentials は認証済みユーザーIDからCouchDB用の認証情報を生成するのだ
func (s *couchDBService) GenerateProxyCredentials(userID string) (string, string, error) {
	// ★修正ポイント: DB検索は不要！
	// PostgreSQLの user_id (例: "16") がそのまま CouchDB の username なのだ
	couchUsername := userID

	// HMAC-SHA1 で署名トークンを作成
	// Hex( HMAC-SHA1(couchUsername, secret) )
	mac := hmac.New(sha1.New, []byte(s.configSecret))
	mac.Write([]byte(couchUsername))
	token := hex.EncodeToString(mac.Sum(nil))

	return couchUsername, token, nil
}

func (s *couchDBService) GetCouchDBURL() string {
	return s.configURL
}
