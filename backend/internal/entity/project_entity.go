package entity

import "time"

type Project struct {
	ProjectID     string    `json:"project_id" gorm:"primaryKey;column:project_id;type:text;default:gen_random_uuid()"`
	ProjectName   string    `json:"project_name" gorm:"column:project_name"`
	Description   string    `json:"disscription" gorm:"column:disscription"`
	StartDay      time.Time `json:"start_day" gorm:"column:start_day;type:date"`
	FinishedDay   time.Time `json:"finished_day" gorm:"column:finished_day;type:date"`
	UpdatedDay    time.Time `json:"updated_day" gorm:"column:updated_day;type:date"`
	Note          string    `json:"note" gorm:"column:note"`
	WorkstationID int64     `json:"workstation_id" gorm:"column:workstation_id"`
	UserID        int64     `json:"user_id" gorm:"column:user_id"`
}

func (Project) TableName() string {
	return "projects"
}
