package entity

import "time"

type User struct {
	UserID      int64     `json:"user_id" gorm:"primaryKey;column:user_id"`
	UserName    string    `json:"user_name" gorm:"column:user_name"`
	DisplayName string    `json:"display_name" gorm:"column:display_name"`
	MailAddress string    `json:"mail_address" gorm:"column:mail_address;unique"`
	Password    string    `json:"-" gorm:"column:password"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	Timezone    string    `json:"timezone" gorm:"column:timezone"`
}

func (User) TableName() string {
	return "users"
}
