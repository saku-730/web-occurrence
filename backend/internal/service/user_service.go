package service

import (
	"errors"

	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
	"strings"

	"gorm.io/gorm" 
)

// (今回追加) サービス層でもカスタムエラーを定義（Handlerが依存するのはServiceのエラー）
var ErrEmailConflict = errors.New("このメールアドレスは既に使用されています")

// UserService はビジネスロジックのインターフェースなのだ
type UserService interface {
	RegisterUser(req *model.UserRegisterRequest) (*entity.User, error)
	LoginUser(req *model.UserLoginRequest) (string, error)
}

// userService は UserService の実装なのだ
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService は userService のインスタンスを生成するのだ
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterUser はユーザー登録のロジックを実行するのだ
func (s *userService) RegisterUser(req *model.UserRegisterRequest) (*entity.User, error) {
	hashedPassword, err := infrastructure.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(req.MailAddress, "@")
	username := parts[0]
	newUser := &entity.User{
		UserName:    username,
		DisplayName: username,
		MailAddress: req.MailAddress,
		Password:    hashedPassword,
	}

	// 4. Repository を呼び出してDBに保存
	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		// --- (ここから変更) ---
		// リポジトリから返されたエラーが「メール重複エラー」かチェック
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			// Handler層に、サービス層で定義した「競合エラー」を返す
			return nil, ErrEmailConflict
		}
		// --- (ここまで変更) ---

		// その他のエラー
		return nil, err
	}
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

	token, err := infrastructure.GenerateToken(user.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}
