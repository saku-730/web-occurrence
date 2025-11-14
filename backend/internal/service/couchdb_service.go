package service

import (
	"errors"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
)

// CouchDBService はCouchDBセッション発行のロジックを担当するのだ
type CouchDBService interface {
	RequestCouchDBSession(userID string) (string, error)
}

type couchDBService struct {
	userRepo  repository.UserRepository
	couchClient infrastructure.CouchDBClient
}

// NewCouchDBService は CouchDBService のインスタンスを生成するのだ
func NewCouchDBService(
	userRepo repository.UserRepository, 
	couchClient infrastructure.CouchDBClient,
) CouchDBService {
	return &couchDBService{
		userRepo:  userRepo,
		couchClient: couchClient,
	}
}

// RequestCouchDBSession は、sのフローのステップ1,2,3を統括するのだ
func (s *couchDBService) RequestCouchDBSession(userID string) (string, error) {

	// 1. JWTのUserID (UUID) から PostgreSQL のユーザー情報を検索
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		// (gorm.ErrRecordNotFound も含む)
		return "", errors.New("JWTに対応するユーザーが見つかりません")
	}

	// 2. CouchDBが認証に使う「username」（例: "test"）を取得
	// (sのRegisterロジックに基づき、UserName を使っている)
	couchDBUsername := user.UserName
	if couchDBUsername == "" {
		return "", errors.New("CouchDBのユーザー名が空です")
	}

	// 3. CouchDBクライアントを呼び出し、代理発行を依頼
	cookieString, err := s.couchClient.GetSessionCookie(couchDBUsername)
	if err != nil {
		return "", err
	}

	// 4. 取得したCookie文字列をHandlerに返す
	return cookieString, nil
}
