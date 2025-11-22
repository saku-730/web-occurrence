package entity

import "time"

type ProjectMember struct {
	ProjectMemberID string    `json:"project_member_id" gorm:"primaryKey;column:project_member_id;type:text;default:gen_random_uuid()"`
	ProjectID       string    `json:"project_id" gorm:"column:project_id;type:text"`
	UserID          int64     `json:"user_id" gorm:"column:user_id"`
	JoinDay         time.Time `json:"join_day" gorm:"column:join_day;type:date"`
	FinishDay       time.Time `json:"finish_day" gorm:"column:finish_day;type:date"`
	WorkstationID   int64     `json:"workstation_id" gorm:"column:workstation_id"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}
