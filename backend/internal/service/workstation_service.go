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
	EnsureAllDatabases() error
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

// CreateWorkstation はワークステーションを作成し、CouchDBのDBを即時作成するのだ
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

	// ★修正: initMasterData の呼び出しを完全に削除したのだ！
	// DB作成と初期データ投入が不要になったので、この処理は不要なのだ。
    
    // 代わりに、ワークステーション作成時に即座にCouchDBのDB（箱）だけを作成する
    dbName := s.couchClient.CreateWorkstationDBName(createdWS.WorkstationID)
    if err := s.couchClient.CreateDatabase(dbName); err != nil {
        // DB作成自体が管理者権限のエラーで失敗した場合は、警告ログを出す
        // データ損失を防ぐため、Postgres登録自体は失敗としない
        fmt.Printf("Warning: Failed to instantly create CouchDB for new WS (%s). Replication will fail until fixed: %v\n", dbName, err)
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

// EnsureAllDatabases はすべてのワークステーションのDBを確保するのだ
func (s *workstationService) EnsureAllDatabases() error {
	workstations, err := s.wsRepo.GetAllWorkstations()
	if err != nil {
		return err
	}

	fmt.Println("--- Checking CouchDB Databases for Workstations ---")
	for _, ws := range workstations {
		dbName := s.couchClient.CreateWorkstationDBName(ws.WorkstationID)
		if err := s.couchClient.CreateDatabase(dbName); err != nil {
			// DB作成エラーは致命的なのでログに残す
			fmt.Printf("Error creating DB %s: %v\n", dbName, err)
		} else {
			fmt.Printf("Database ensured: %s\n", dbName)
		}
	}
	fmt.Println("-------------------------------------------------")
	return nil
}
