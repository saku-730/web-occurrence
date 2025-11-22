package repository

import (
	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"gorm.io/gorm"
)

type WorkstationRepository interface {
	CreateWorkstation(ws *entity.Workstation) (*entity.Workstation, error)
	AddUserToWorkstation(userID, workstationID int64, roleID int64) error
}

type workstationRepository struct {
	db *gorm.DB
}

func NewWorkstationRepository(db *gorm.DB) WorkstationRepository {
	return &workstationRepository{db: db}
}

func (r *workstationRepository) CreateWorkstation(ws *entity.Workstation) (*entity.Workstation, error) {
	result := r.db.Create(ws)
	if result.Error != nil {
		return nil, result.Error
	}
	return ws, nil
}

func (r *workstationRepository) AddUserToWorkstation(userID, workstationID int64, roleID int64) error {
	link := entity.WorkstationUser{
		UserID:        userID,
		WorkstationID: workstationID,
		RoleID:        roleID, // 1: admin (仮)
	}
	return r.db.Create(&link).Error
}
