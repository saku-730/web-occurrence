package entity

import "time"

// User は users テーブルのカラムに対応するのだ
// GORM用の struct tags (`gorm:"..."`) を追加したのだ
type User struct {
	// PostgreSQL側で pgcrypto の gen_random_uuid() を使う想定なのだ
	UserID       string    `json:"user_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName     string    `json:"user_name" gorm:"column:user_name"`
	DisplayName  string    `json:"display_name" gorm:"column:display_name"`
	MailAddress  string    `json:"mail_address" gorm:"column:mail_address;unique"` // メールアドレスはユニーク制約
	Password     string    `json:"-" gorm:"column:password"`                      // JSONで返さない
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`  // GORMが自動で作成時刻を入れる
	Timezone     string    `json:"timezone" gorm:"column:timezone"`
}

// TableName メソッドは、GORMに構造体とテーブル名をマッピングさせるのだ
func (User) TableName() string {
	return "users"
}
