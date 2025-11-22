package repository

import (
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"gorm.io/gorm"
)

type MasterRepository interface {
	GetAllLanguages() ([]model.Language, error)
	GetAllFileTypes() ([]model.FileType, error)
	GetAllFileExtensions() ([]model.FileExtension, error)
	GetAllUserRoles() ([]model.UserRole, error)
	// ▼ 変更: 全件取得をやめて、ワークステーション指定で取得するメソッドにするのだ
	GetUsersByWorkstationID(workstationID int64) ([]model.WorkstationUser, error)
}

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepository{db: db}
}

func (r *masterRepository) GetAllLanguages() ([]model.Language, error) {
	var list []model.Language
	err := r.db.Find(&list).Error
	return list, err
}

func (r *masterRepository) GetAllFileTypes() ([]model.FileType, error) {
	var list []model.FileType
	err := r.db.Find(&list).Error
	return list, err
}

func (r *masterRepository) GetAllFileExtensions() ([]model.FileExtension, error) {
	var list []model.FileExtension
	err := r.db.Find(&list).Error
	return list, err
}

func (r *masterRepository) GetAllUserRoles() ([]model.UserRole, error) {
	var list []model.UserRole
	err := r.db.Find(&list).Error
	return list, err
}

// ▼ 追加実装: 指定されたワークステーションに所属するユーザーのみを取得
func (r *masterRepository) GetUsersByWorkstationID(workstationID int64) ([]model.WorkstationUser, error) {
	var list []model.WorkstationUser
	
	// workstation_user テーブルと users テーブルを JOIN して、
	// そのワークステーションに紐付いているユーザーだけを抽出するのだ
	err := r.db.Table("users").
		Select("users.user_id, users.display_name").
		Joins("JOIN workstation_user ON users.user_id = workstation_user.user_id").
		Where("workstation_user.workstation_id = ?", workstationID).
		Find(&list).Error
		
	return list, err
}
