package entity

type SpecimenMethod struct {
	SpecimenMethodsID string `json:"specimen_methods_id" gorm:"primaryKey;column:specimen_methods_id;type:text;default:gen_random_uuid()"`
	MethodCommonName  string `json:"method_common_name" gorm:"column:method_common_name"`
	PageID            string `json:"page_id" gorm:"column:page_id;type:text"`
	WorkstationID     int64  `json:"workstation_id" gorm:"column:workstation_id"`
	UserID            int64  `json:"user_id" gorm:"column:user_id"`
}

func (SpecimenMethod) TableName() string {
	return "specimen_methods"
}
