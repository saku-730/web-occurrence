package entity

import "time"

type Observation struct {
	ObservationID       string    `json:"observation_id" gorm:"primaryKey;column:observation_id;type:text"`
	OccurrenceID        string    `json:"occurrence_id" gorm:"column:occurrence_id;type:text"`
	UserID              int64     `json:"user_id" gorm:"column:user_id"` // Text -> BigInt
	ObservationMethodID string    `json:"observation_method_id" gorm:"column:observation_method_id;type:text"`
	Behavior            string    `json:"behavior" gorm:"column:behavior"`
	ObservedAt          time.Time `json:"observed_at" gorm:"column:observed_at"`
}

func (Observation) TableName() string {
	return "observations"
}
