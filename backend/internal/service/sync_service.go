package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"gorm.io/gorm"
)

type SyncService interface {
	StartPolling()
}

type syncService struct {
	db          *gorm.DB
	couchClient infrastructure.CouchDBClient
}

func NewSyncService(db *gorm.DB, couchClient infrastructure.CouchDBClient) SyncService {
	return &syncService{
		db:          db,
		couchClient: couchClient,
	}
}

// 簡易的なポーリングによる同期（本番では _changes フィードのロングポーリング推奨）
func (s *syncService) StartPolling() {
	go func() {
		ticker := time.NewTicker(10 * time.Second) // 10秒ごとにチェック
		for range ticker.C {
			s.syncAll()
		}
	}()
}

// ここでは「CouchDBにある全てのoccurrenceドキュメント」を取得してUpsertする単純な実装にするのだ
// (差分更新は _changes を使うともっと効率的になるのだ)
func (s *syncService) syncAll() {
	// CouchDBから全ドキュメントIDを取得するメソッドが必要だけど、
	// ここでは解説用に「1件処理するロジック」を書くのだ。
	// 実際には couchClient.GetAllDocs() などを実装してループさせるのだ。
}

// 1つのドキュメント（JSON）をPostgresに保存するメインロジック
func (s *syncService) ProcessDocument(doc map[string]interface{}) error {
	// typeチェック
	docType, _ := doc["type"].(string)
	if docType != "occurrence" {
		return nil // 対象外
	}

	// JSONのパース
	// 実際には map[string]interface{} からキャストしまくるのは大変なので、
	// 一度 JSON bytes にして struct にマッピングすると楽なのだ。
	jsonBytes, _ := json.Marshal(doc)
	
	var data IncomingOccurrenceData
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Classification
		if data.ClassificationData.ClassificationID != "" {
			cls := entity.ClassificationJSON{
				ClassificationID:    data.ClassificationData.ClassificationID,
				ClassClassification: data.ClassificationData.ClassClassification,
			}
			if err := tx.Save(&cls).Error; err != nil { return err }
		}

		// 2. Place
		if data.PlaceData.PlaceID != "" {
			pl := entity.Place{
				PlaceID:     data.PlaceData.PlaceID,
				PlaceNameID: nil, // 省略
				Coordinates: data.PlaceData.Coordinates,
				Accuracy:    data.PlaceData.Accuracy,
			}
			if err := tx.Save(&pl).Error; err != nil { return err }
		}

		// 3. Occurrence (本体)
		// IDから "occ_" などのプレフィックスを取り除くかどうかは要件次第だけど
		// UUIDとして保存するならそのまま入れるか、整形するのだ。
		// 今回はそのまま文字列として入れるのだ。
		
		// 時刻パース
		createdAt, _ := time.Parse(time.RFC3339, data.CreatedAt)

		occ := entity.Occurrence{
			OccurrenceID:     data.ID, // _id を使う
			WorkstationID:    data.WorkstationID,
			UserID:           data.CreatedByUserID,
			ProjectID:        data.ProjectID,
			IndividualID:     data.OccurrenceData.IndividualID,
			Lifestage:        data.OccurrenceData.Lifestage,
			Sex:              data.OccurrenceData.Sex,
			BodyLength:       data.OccurrenceData.BodyLength,
			Note:             data.OccurrenceData.Note,
			ClassificationID: data.ClassificationData.ClassificationID,
			PlaceID:          data.PlaceData.PlaceID,
			LanguageID:       data.LanguageID,
			CreatedAt:        createdAt,
			Timezone:         data.Timezone,
		}
		if err := tx.Save(&occ).Error; err != nil { return err }

		// 4. Related Tables (Identifications, etc.)
		// ここでは省略するけど、同様にループして Save するのだ

		log.Printf("Synced occurrence: %s", occ.OccurrenceID)
		return nil
	})
}

// JSON受け取り用の構造体定義
type IncomingOccurrenceData struct {
	ID              string `json:"_id"`
	WorkstationID   string `json:"workstation_id"`
	CreatedByUserID string `json:"created_by_user_id"`
	ProjectID       *string `json:"project_id"`
	CreatedAt       string  `json:"created_at"`
	Timezone        string  `json:"timezone"`
	LanguageID      *string `json:"language_id"`

	OccurrenceData struct {
		IndividualID string   `json:"individual_id"`
		Lifestage    string   `json:"lifestage"`
		Sex          string   `json:"sex"`
		BodyLength   *float64 `json:"body_length"`
		Note         string   `json:"note"`
	} `json:"occurrence_data"`

	ClassificationData struct {
		ClassificationID    string                 `json:"classification_id"`
		ClassClassification map[string]interface{} `json:"class_classification"`
	} `json:"classification_data"`

	PlaceData struct {
		PlaceID     string                 `json:"place_id"`
		Coordinates map[string]interface{} `json:"coordinates"`
		Accuracy    *float64               `json:"accuracy"`
	} `json:"place_data"`
}
