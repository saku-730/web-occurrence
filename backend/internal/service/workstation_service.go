package service

import (
	"fmt"
	"strconv"

	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
)

type WorkstationService interface {
	CreateWorkstation(userID string, req *model.CreateWorkstationRequest) (*entity.Workstation, error)
	GetMyWorkstations(userID string) ([]entity.Workstation, error)
}

type workstationService struct {
	wsRepo      repository.WorkstationRepository
	masterRepo  repository.MasterRepository
	couchClient infrastructure.CouchDBClient
}

func NewWorkstationService(
	wsRepo repository.WorkstationRepository,
	masterRepo repository.MasterRepository,
	couchClient infrastructure.CouchDBClient,
) WorkstationService {
	return &workstationService{
		wsRepo:      wsRepo,
		masterRepo:  masterRepo,
		couchClient: couchClient,
	}
}

func (s *workstationService) CreateWorkstation(userIDStr string, req *model.CreateWorkstationRequest) (*entity.Workstation, error) {
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	newWS := &entity.Workstation{
		WorkstationName: req.WorkstationName,
	}
	createdWS, err := s.wsRepo.CreateWorkstation(newWS)
	if err != nil {
		return nil, err
	}

	if err := s.wsRepo.AddUserToWorkstation(userID, createdWS.WorkstationID, 1); err != nil {
		return nil, err
	}

	languages, _ := s.masterRepo.GetAllLanguages()
	fileTypes, _ := s.masterRepo.GetAllFileTypes()
	fileExts, _ := s.masterRepo.GetAllFileExtensions()
	roles, _ := s.masterRepo.GetAllUserRoles()
	
	// 修正: entity.WorkstationUser型を使わず、CouchDB保存用のMapスライスにする
	users := []map[string]interface{}{
		{
			"user_id":      userID,
			"display_name": "Current User", // ※本来はDBからユーザー名を取得すべきだけど、今は固定値で回避
		},
	}

	docID := "_local/master_data" 
	docData := map[string]interface{}{
		"_id":                docID,
		"type":               "master_data",
		"workstation_id":     fmt.Sprintf("%d", createdWS.WorkstationID), 
		"data": map[string]interface{}{
			"languages":         languages,
			"file_types":        fileTypes,
			"file_extensions":   fileExts,
			"user_roles":        roles,
			"workstation_users": users,
		},
	}

	if err := s.couchClient.UpsertDocument(docID, docData); err != nil {
		fmt.Printf("Failed to init CouchDB master data: %v\n", err)
	}

	return createdWS, nil
}

func (s *workstationService) GetMyWorkstations(userIDStr string) ([]entity.Workstation, error) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.wsRepo.GetWorkstationsByUserID(userID)
}
