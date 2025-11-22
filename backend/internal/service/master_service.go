package service

import (
	"strconv"

	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
)

// データをまとめるための構造体
type MasterDataResponse struct {
	WorkstationID    int64                   `json:"workstation_id"`
	Languages        []model.Language        `json:"languages"`
	FileTypes        []model.FileType        `json:"file_types"`
	FileExtensions   []model.FileExtension   `json:"file_extensions"`
	UserRoles        []model.UserRole        `json:"user_roles"`
	WorkstationUsers []model.WorkstationUser `json:"workstation_users"`
}

type MasterService interface {
	GetMasterData(userID string) (*MasterDataResponse, error)
}

type masterService struct {
	masterRepo repository.MasterRepository
	wsRepo     repository.WorkstationRepository
}

func NewMasterService(masterRepo repository.MasterRepository, wsRepo repository.WorkstationRepository) MasterService {
	return &masterService{
		masterRepo: masterRepo,
		wsRepo:     wsRepo,
	}
}

func (s *masterService) GetMasterData(userIDStr string) (*MasterDataResponse, error) {
	// 1. まずワークステーションIDを特定するのだ
	var wsID int64 = 0
	if userIDStr != "" {
		uid, _ := strconv.ParseInt(userIDStr, 10, 64)
		// ユーザーが所属しているワークステーションを取得
		ws, err := s.wsRepo.FindWorkstationByUserID(uid)
		if err == nil && ws != nil {
			wsID = ws.WorkstationID
		}
	}

	// 2. 共通のマスターデータを取得するのだ
	languages, err := s.masterRepo.GetAllLanguages()
	if err != nil { return nil, err }
	
	fileTypes, err := s.masterRepo.GetAllFileTypes()
	if err != nil { return nil, err }
	
	fileExts, err := s.masterRepo.GetAllFileExtensions()
	if err != nil { return nil, err }
	
	roles, err := s.masterRepo.GetAllUserRoles()
	if err != nil { return nil, err }
	
	// 3. ▼ 修正: ユーザー一覧は、特定した wsID に紐づくものだけを取得するのだ
	var users []model.WorkstationUser
	if wsID != 0 {
		users, err = s.masterRepo.GetUsersByWorkstationID(wsID)
		if err != nil { return nil, err }
	} else {
		// ワークステーションに所属していない場合は空リストにするのだ
		users = []model.WorkstationUser{}
	}

	return &MasterDataResponse{
		WorkstationID:    wsID,
		Languages:        languages,
		FileTypes:        fileTypes,
		FileExtensions:   fileExts,
		UserRoles:        roles,
		WorkstationUsers: users,
	}, nil
}
