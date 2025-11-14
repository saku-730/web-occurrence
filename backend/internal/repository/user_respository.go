package repository

import (
	"errors"
	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"strings"

	"gorm.io/gorm"
)

// (今回追加) リポジトリ層で発生した特異的なエラーを定義
var ErrEmailAlreadyExists = errors.New("このメールアドレスは既に使用されています")

// UserRepository はDB操作のインターフェースなのだ
type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
}

// userRepository は UserRepository の実装なのだ
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository は userRepository のインスタンスを生成するのだ
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entity.User) (*entity.User, error) {
	user.Timezone = "9"
	result := r.db.Create(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			// サービス層に「メール重複エラー」であることを伝える
			return nil, ErrEmailAlreadyExists
		}

		errorString := result.Error.Error()
		if strings.Contains(errorString, "23505") {
			return nil, ErrEmailAlreadyExists
		}




		// GORMが提供する「キー重複エラー」かどうかを errors.Is でチェック
		// これなら "23505" みたいなドライバ固有コードを知らなくていいのだ
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailAlreadyExists
		}

		// その他のDBエラー
		return nil, result.Error
	}
	return user, nil
}


// FindUserByEmail はメールアドレスでユーザーを1件検索するのだ
func (r *userRepository) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("mail_address = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
