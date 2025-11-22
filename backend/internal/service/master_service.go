package service

import (
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
)

// データをまとめるための構造体
type MasterDataResponse struct {
	Languages        []model.Language        `json:"languages"`
	FileTypes        []model.FileType        `json:"file_types"`
	FileExtensions   []model.FileExtension   `json:"file_extensions"`
	UserRoles        []model.UserRole        `json:"user_roles"`
	WorkstationUsers []model.WorkstationUser `json:"workstation_users"`
}

type MasterService interface {
	GetMasterData() (*MasterDataResponse, error)
}

type masterService struct {
	masterRepo repository.MasterRepository
}

func NewMasterService(repo repository.MasterRepository) MasterService {
	return &masterService{masterRepo: repo}
}

func (s *masterService) GetMasterData() (*MasterDataResponse, error) {
	languages, err := s.masterRepo.GetAllLanguages()
	if err != nil { return nil, err }
	
	fileTypes, err := s.masterRepo.GetAllFileTypes()
	if err != nil { return nil, err }
	
	fileExts, err := s.masterRepo.GetAllFileExtensions()
	if err != nil { return nil, err }
	
	roles, err := s.masterRepo.GetAllUserRoles()
	if err != nil { return nil, err }
	
	users, err := s.masterRepo.GetAllWorkstationUsers()
	if err != nil { return nil, err }

	return &MasterDataResponse{
		Languages:        languages,
		FileTypes:        fileTypes,
		FileExtensions:   fileExts,
		UserRoles:        roles,
		WorkstationUsers: users,
	}, nil
}
