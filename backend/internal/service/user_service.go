package service

import (
	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
	"strings"
)

type UserService interface {
	RegisterUser(req *model.UserRegisterRequest) (*entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterUser はユーザー登録のロジックを実行するのだ
func (s *userService) RegisterUser(req *model.UserRegisterRequest) (*entity.User, error) {
	
	// 1. パスワードをハッシュ化する (infrastructure を使う)
	hashedPassword, err := infrastructure.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 2. メールアドレスから @ より前の部分を切り出す
	parts := strings.Split(req.MailAddress, "@")
	username := parts[0]

	// 3. Repositoryに渡すための entity を作成
	newUser := &entity.User{
		UserName:    username,
		DisplayName: username,
		MailAddress: req.MailAddress,
		Password:    hashedPassword,
	}

	// 4. Repository を呼び出してDBに保存
	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		// (本当はここでメールアドレス重複エラー(23505)などをハンドリングすべき)
		return nil, err
	}

	return createdUser, nil
}
