
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

    // ★追加: DB作成後、即座にユーザーにアクセス権を付与するのだ！ (403対策)
    if err := s.couchClient.SetDatabaseUserAccess(dbName, userIDStr); err != nil {
        // アクセス権限設定に失敗したら、同期ができなくなるのでエラーを返すのだ
        return nil, fmt.Errorf("ワークステーションは作成されましたが、CouchDBのアクセス権設定に失敗しました: %w", err)
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

func (s *workstationService) EnsureAllDatabases() error {
    // ★注意: wsRepoにGetAllWorkstationUserRelations()があることを前提とするのだ。
    // WorkstationUserエンティティのリストを取得し、ユーザーIDとWS IDのペアを扱うのだ。
	
    // 仮定: リポジトリが全てのユーザー・ワークステーションの紐付けを返すのだ
    relations, err := s.wsRepo.GetAllWorkstationUserRelations() // ★リポジトリ側で実装が必要なのだ
	if err != nil {
		return err
	}

	fmt.Println("--- Checking CouchDB Databases for Workstations ---")
    // ★ユーザーIDとワークステーションIDのペアごとにDBを作成/確認するのだ
	for _, rel := range relations {
        // userIDStr: ユーザーIDを文字列に変換するのだ
        userIDStr := strconv.FormatInt(rel.UserID, 10)
        
        // dbNameを新しい命名規則で生成するのだ
		dbName := fmt.Sprintf("%s_db_ws_%d", userIDStr, rel.WorkstationID)
        
        // DB作成
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
