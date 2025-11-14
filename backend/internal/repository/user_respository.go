package repository

import (
	// "database/sql" // 標準ライブラリは不要になったのだ
	"myapp/entity" // (sのプロジェクト名)

	"gorm.io/gorm" // GORMをインポート
)

// UserRepository はDB操作のインターフェースなのだ
type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	// FindUserByEmail(email string) (*entity.User, error) // (ログイン機能で必要になる)
}

// userRepository は UserRepository の実装なのだ
type userRepository struct {
	db *gorm.DB // *sql.DB から *gorm.DB に変更
}

// NewUserRepository は userRepository のインスタンスを生成するのだ
// GORMのDBインスタンスを受け取るように変更
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser は users テーブルに新しいユーザーを挿入するのだ (GORM版)
func (r *userRepository) CreateUser(user *entity.User) (*entity.User, error) {
	// タイムゾーンはService層から渡されたもの（またはここで設定）
	user.Timezone = "UTC" 

	// GORMの Create メソッドを使うのだ
	// userポインタを渡すと、GORMが自動でSQLのINSERT文を生成する。
	// entityで設定した gorm:"default:gen_random_uuid()" や gorm:"autoCreateTime" が
	// DB側で評価され、GORMがその結果（UserIDやCreatedAt）を
	// user ポインタの中身にマッピングし直してくれるのだ。
	result := r.db.Create(user)

	if result.Error != nil {
		// エラー（メールアドレスの重複(unique制約違反)など）
		return nil, result.Error
	}

	// user にはDBから返された UserID と CreatedAt が入っている
	return user, nil
}
