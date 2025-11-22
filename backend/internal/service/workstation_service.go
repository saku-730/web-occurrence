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

	// 1. PostgreSQLにワークステーションを作成
	newWS := &entity.Workstation{
		WorkstationName: req.WorkstationName,
	}
	createdWS, err := s.wsRepo.CreateWorkstation(newWS)
	if err != nil {
		return nil, err
	}

	// 2. 作成者を管理者として紐付け (RoleID=1: adminと仮定)
	if err := s.wsRepo.AddUserToWorkstation(userID, createdWS.WorkstationID, 1); err != nil {
		return nil, err
	}

	// 3. CouchDBにマスターデータを投入
	// 本来はワークステーションごとにデータを分けるかIDで区別するが、
	// 今回は単一DBで _local/master_data を使う要件なので、それを更新する。
	// ※ 複数WSがある場合、_local/master_data をどう共有するかは要検討だけど、
	//    一旦は「作成時にマスターデータを初期化/更新する」という挙動にするのだ。

	languages, _ := s.masterRepo.GetAllLanguages()
	fileTypes, _ := s.masterRepo.GetAllFileTypes()
	fileExts, _ := s.masterRepo.GetAllFileExtensions()
	roles, _ := s.masterRepo.GetAllUserRoles()
	
	// 所属ユーザー一覧を取得（今は作成者だけ）
	// 本当は wsRepo.GetUsersByWorkstationID みたいなメソッドが必要だけど簡易的に作成
	users := []entity.WorkstationUser{
		{UserID: userID, DisplayName: "Current User"}, // 名前はUserRepoから引くべき
	}

	docID := "_local/master_data" 
	docData := map[string]interface{}{
		"_id":                docID,
		"type":               "master_data",
		// どのワークステーションのデータか区別するためのIDを入れるのもあり
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
		// ログ出すだけにしておく
		fmt.Printf("Failed to init CouchDB master data: %v\n", err)
	}

	return createdWS, nil
}
