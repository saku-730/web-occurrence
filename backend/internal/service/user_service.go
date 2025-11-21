package service

import (
	"errors"
	"fmt"
	"strconv" // 追加: int64をstringにするため
	"strings"

	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
	"gorm.io/gorm"
)

var ErrEmailConflict = errors.New("このメールアドレスは既に使用されています")
var ErrCouchDBUserCreation = errors.New("CouchDBユーザーの作成に失敗しました")

type UserService interface {
	RegisterUser(req *model.UserRegisterRequest) (*entity.User, error)
	LoginUser(req *model.UserLoginRequest) (string, error)
}

type userService struct {
	userRepo    repository.UserRepository
	couchClient infrastructure.CouchDBClient
}

func NewUserService(
	userRepo repository.UserRepository,
	couchClient infrastructure.CouchDBClient,
) UserService {
	return &userService{
		userRepo:    userRepo,
		couchClient: couchClient,
	}
}

// RegisterUser はユーザー登録のロジックを実行するのだ
func (s *userService) RegisterUser(req *model.UserRegisterRequest) (*entity.User, error) {

	// 1. ハッシュパスワード生成
	hashedPassword, err := infrastructure.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 2. PostgreSQL用エンティティ作成
	parts := strings.Split(req.MailAddress, "@")
	username := parts[0]
	newUser := &entity.User{
		UserName:    username,
		DisplayName: username,
		MailAddress: req.MailAddress,
		Password:    hashedPassword,
	}

	// 3. PostgreSQLに保存 (ここで UserID が発番される！)
	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, ErrEmailConflict
		}
		return nil, err
	}

	// 4. CouchDBユーザー作成
	// ★修正ポイント: PostgreSQLの user_id を文字列にして、CouchDBのユーザー名として使うのだ
	couchDBUsername := strconv.FormatInt(createdUser.UserID, 10)
	
	// パスワードは平文で渡す（CouchDB側でハッシュ化される）
	err = s.couchClient.CreateCouchDBUser(couchDBUsername, req.Password)
	if err != nil {
		// 失敗したら本当はロールバックしたいけど、今はエラーを返す
		return nil, fmt.Errorf("%w: %v", ErrCouchDBUserCreation, err)
	}

	return createdUser, nil
}

func (s *userService) LoginUser(req *model.UserLoginRequest) (string, error) {
	user, err := s.userRepo.FindUserByEmail(req.MailAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("メールアドレスまたはパスワードが正しくありません")
		}
		return "", err
	}

	isValidPassword := infrastructure.CheckPasswordHash(req.Password, user.Password)
	if !isValidPassword {
		return "", errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	userIDString := strconv.FormatInt(user.UserID, 10)
	token, err := infrastructure.GenerateToken(userIDString)
	if err != nil {
		return "", err
	}

	return token, nil
}
