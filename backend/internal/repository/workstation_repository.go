package repository

import (
	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"gorm.io/gorm"
)

type WorkstationRepository interface {
	CreateWorkstation(ws *entity.Workstation) (*entity.Workstation, error)
	AddUserToWorkstation(userID, workstationID int64, roleID int64) error
	FindWorkstationByUserID(userID int64) (*entity.Workstation, error)
	GetWorkstationsByUserID(userID int64) ([]entity.Workstation, error)
	GetAllWorkstations() ([]entity.Workstation, error)
	GetAllWorkstationUserRelations() ([]entity.WorkstationUser, error)
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
		RoleID:        int(roleID),
	}
	return r.db.Create(&link).Error
}

func (r *workstationRepository) FindWorkstationByUserID(userID int64) (*entity.Workstation, error) {
	var ws entity.Workstation
	err := r.db.Table("workstation").
		Joins("JOIN workstation_user ON workstation.workstation_id = workstation_user.workstation_id").
		Where("workstation_user.user_id = ?", userID).
		First(&ws).Error

	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (r *workstationRepository) GetWorkstationsByUserID(userID int64) ([]entity.Workstation, error) {
	var workstations []entity.Workstation
	err := r.db.Table("workstation").
		Joins("JOIN workstation_user ON workstation.workstation_id = workstation_user.workstation_id").
		Where("workstation_user.user_id = ?", userID).
		Find(&workstations).Error

	if err != nil {
		return nil, err
	}
	return workstations, nil
}

// ▼ 追加実装: 全ワークステーション取得（SyncServiceで使用）
func (r *workstationRepository) GetAllWorkstations() ([]entity.Workstation, error) {
	var workstations []entity.Workstation
	if err := r.db.Find(&workstations).Error; err != nil {
		return nil, err
	}
	return workstations, nil
}

// ★新規実装: 全てのユーザー・ワークステーションの紐付けを取得するのだ
func (r *workstationRepository) GetAllWorkstationUserRelations() ([]entity.WorkstationUser, error) {
	var relations []entity.WorkstationUser
	// entity.WorkstationUser は GORM にテーブル名 (workstation_user) を教えるのだ
	if err := r.db.Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}
