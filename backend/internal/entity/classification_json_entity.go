package entity

type ClassificationJSON struct {
	ClassificationID    string `json:"classification_id" gorm:"primaryKey;column:classification_id;type:text"`
	ClassClassification string `json:"class_classification" gorm:"column:class_classification;type:jsonb"`
}

func (ClassificationJSON) TableName() string {
	return "classification_json"
}
