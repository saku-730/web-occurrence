package entity

import "time"

type MakeSpecimen struct {
	MakeSpecimenID string    `json:"make_specimen_id" gorm:"primaryKey;column:make_specimen_id;type:text"`
	SpecimenID     string    `json:"specimen_id" gorm:"column:specimen_id;type:text"`
	UserID         int64     `json:"user_id" gorm:"column:user_id"` // Text -> BigInt
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
}

func (MakeSpecimen) TableName() string {
	return "make_specimen"
}
