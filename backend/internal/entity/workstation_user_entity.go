package entity

type WorkstationUser struct {
	WorkstationID int64 `json:"workstation_id" gorm:"primaryKey;column:workstation_id"`
	UserID        int64 `json:"user_id" gorm:"primaryKey;column:user_id"`
	RoleID        int   `json:"role_id" gorm:"column:role_id;default:1"`
}

func (WorkstationUser) TableName() string {
	return "workstation_user"
}
