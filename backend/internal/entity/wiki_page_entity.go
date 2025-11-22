package entity

import "time"

type WikiPage struct {
	PageID        string    `json:"page_id" gorm:"primaryKey;column:page_id;type:text;default:gen_random_uuid()"`
	Title         string    `json:"title" gorm:"column:title"`
	UserID        int64     `json:"user_id" gorm:"column:user_id"`
	CreatedDate   time.Time `json:"created_date" gorm:"column:created_date;autoCreateTime"`
	UpdatedDate   time.Time `json:"updated_date" gorm:"column:updated_date"`
	ContentPath   string    `json:"content_path" gorm:"column:content_path"`
	WorkstationID int64     `json:"workstation_id" gorm:"column:workstation_id"`
}

func (WikiPage) TableName() string {
	return "wiki_pages"
}
