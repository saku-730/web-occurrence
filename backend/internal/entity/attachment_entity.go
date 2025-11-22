package entity

type Attachment struct {
	AttachmentID string `json:"attachment_id" gorm:"primaryKey;column:attachment_id;type:text"`
	FilePath     string `json:"file_path" gorm:"column:file_path"`
	UserID       int64  `json:"user_id" gorm:"column:user_id"` // Text -> BigInt
}

func (Attachment) TableName() string {
	return "attachments"
}
