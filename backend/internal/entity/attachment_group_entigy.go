package entity

type AttachmentGroup struct {
	OccurrenceID string `json:"occurrence_id" gorm:"primaryKey;column:occurrence_id;type:text"`
	AttachmentID string `json:"attachment_id" gorm:"primaryKey;column:attachment_id;type:text"`
	Priority     int    `json:"priority" gorm:"column:priority"`
	// マイグレーションで一応残したがNULL許容になったカラム
	WorkstationID *int64 `json:"workstation_id" gorm:"column:workstation_id"`
	UserID        *int64 `json:"user_id" gorm:"column:user_id"`
}

func (AttachmentGroup) TableName() string {
	return "attachment_group"
}
