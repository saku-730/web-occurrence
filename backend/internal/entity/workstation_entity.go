package entity

type Workstation struct {
	WorkstationID   int64  `json:"workstation_id" gorm:"primaryKey;column:workstation_id"`
	WorkstationName string `json:"workstation_name" gorm:"column:workstation_name"`
}

func (Workstation) TableName() string {
	return "workstation"
}
