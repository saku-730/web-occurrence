package service

import (
	"errors"
	"fmt"

	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
	"strings"
	"strconv"

	"gorm.io/gorm" 
)


var ErrEmailConflict = errors.New("このメールアドレスは既に使用されています")
var ErrCouchDBUserCreation = errors.New("CouchDBユーザーの作成に失敗しました") // (今回追加)

// UserService はビジネスロジックのインターフェースなのだ
type UserService interface {
	RegisterUser(req *model.UserRegisterRequest) (*entity.User, error)
	LoginUser(req *model.UserLoginRequest) (string, error)
}

// userService は UserService の実装なのだ
type userService struct {
	userRepo    repository.UserRepository
	couchClient infrastructure.CouchDBClient // (今回追加)
}

// NewUserService は userService のインスタンスを生成するのだ
func NewUserService(
	userRepo repository.UserRepository,
	couchClient infrastructure.CouchDBClient, // (今回追加)
) UserService {
	return &userService{
		userRepo:    userRepo,
		couchClient: couchClient, // (今回追加)
	}
}

// RegisterUser はユーザー登録のロジックを実行するのだ
func (s *userService) RegisterUser(req *model.UserRegisterRequest) (*entity.User, error) {
	
	// 1. PostgreSQL用のハッシュパスワードを生成
	hashedPassword, err := infrastructure.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	
	// 2. PostgreSQLに保存するユーザーエンティティを作成
	parts := strings.Split(req.MailAddress, "@")
	username := parts[0]
	newUser := &entity.User{
		UserName:    username, // (CouchDBのnameにもこれを使う)
		DisplayName: username,
		MailAddress: req.MailAddress,
		Password:    hashedPassword, // (PostgreSQL用のハッシュ)
	}

	// 3. PostgreSQLにユーザーを保存
	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, ErrEmailConflict
		}
		// その他のDBエラー
		return nil, err
	}

	// (CouchDBにはハッシュ化する前の「平文パスワード」を渡すのだ)
	err = s.couchClient.CreateCouchDBUser(createdUser.UserName, req.Password)
	if err != nil {
		// CouchDB側の登録に失敗
		// (本当はここでPostgreSQL側の登録をロールバック（削除）する補償トランザクションを入れると完璧なのだ)
		return nil, fmt.Errorf("%w: %v", ErrCouchDBUserCreation, err)
	}

	// 5. 両方のDBに登録成功
	return createdUser, nil
}

// LoginUser はログイン処理を行い、成功したらJWTトークンを返すのだ
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
