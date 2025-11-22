package entity

import "time"

type Identification struct {
	IdentificationID string    `json:"identification_id" gorm:"primaryKey;column:identification_id;type:text"`
	OccurrenceID     string    `json:"occurrence_id" gorm:"column:occurrence_id;type:text"`
	UserID           int64     `json:"user_id" gorm:"column:user_id"` // Text -> BigInt
	SourceInfo       string    `json:"source_info" gorm:"column:source_info"`
	IdentificatedAt  time.Time `json:"identificated_at" gorm:"column:identificated_at"`
}

func (Identification) TableName() string {
	return "identifications"
}
