package entity

type ObservationMethod struct {
	ObservationMethodID string `json:"observation_method_id" gorm:"primaryKey;column:observation_method_id;type:text;default:gen_random_uuid()"`
	MethodCommonName    string `json:"method_common_name" gorm:"column:method_common_name"`
	PageID              string `json:"pageid" gorm:"column:pageid;type:text"` // カラム名ママ pageid
	WorkstationID       int64  `json:"workstation_id" gorm:"column:workstation_id"`
	UserID              int64  `json:"user_id" gorm:"column:user_id"`
}

func (ObservationMethod) TableName() string {
	return "observation_methods"
}
