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
	GetAllWorkstationUsers() ([]model.WorkstationUser, error)
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

func (r *masterRepository) GetAllWorkstationUsers() ([]model.WorkstationUser, error) {
	var list []model.WorkstationUser
	// usersテーブルから必要なカラムだけ取得
	err := r.db.Table("users").Select("user_id, display_name").Find(&list).Error
	return list, err
}
