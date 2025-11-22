package entity

import "time"

type Occurrence struct {
	OccurrenceID     string    `json:"occurrence_id" gorm:"primaryKey;column:occurrence_id;type:text"`
	WorkstationID    int64     `json:"workstation_id" gorm:"column:workstation_id"` // Text -> BigInt
	UserID           int64     `json:"user_id" gorm:"column:user_id"`             // Text -> BigInt
	ProjectID        string    `json:"project_id" gorm:"column:project_id;type:text"`
	IndividualID     string    `json:"individual_id" gorm:"column:individual_id;type:text"`
	Lifestage        string    `json:"lifestage" gorm:"column:lifestage"`
	Sex              string    `json:"sex" gorm:"column:sex"`
	BodyLength       float64   `json:"body_length" gorm:"column:body_length"`
	Note             string    `json:"note" gorm:"column:note"`
	ClassificationID string    `json:"classification_id" gorm:"column:classification_id;type:text"`
	PlaceID          string    `json:"place_id" gorm:"column:place_id;type:text"`
	LanguageID       string    `json:"language_id" gorm:"column:language_id;type:text"` // Integer? 確認要だがschema.sqlではtextだった
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at"`
	Timezone         string    `json:"timezone" gorm:"column:timezone"`
}

func (Occurrence) TableName() string {
	return "occurrence"
}
