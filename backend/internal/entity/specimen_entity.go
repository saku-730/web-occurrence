package entity

type Specimen struct {
	SpecimenID       string `json:"specimen_id" gorm:"primaryKey;column:specimen_id;type:text"`
	OccurrenceID     string `json:"occurrence_id" gorm:"column:occurrence_id;type:text"`
	InstitutionID    string `json:"institution_id" gorm:"column:institution_id;type:text"`
	CollectionID     string `json:"collection_id" gorm:"column:collection_id;type:text"`
	SpecimenMethodID string `json:"specimen_method_id" gorm:"column:specimen_method_id;type:text"`
}

func (Specimen) TableName() string {
	return "specimen"
}
