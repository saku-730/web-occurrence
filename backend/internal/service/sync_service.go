package service

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/saku-730/web-occurrence/backend/internal/entity"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"gorm.io/gorm"
)

type SyncService interface {
	StartPolling()
	ProcessDocument(doc map[string]interface{}) error
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

func (s *syncService) StartPolling() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			s.syncAll()
		}
	}()
}

func (s *syncService) syncAll() {
	// 実装省略
}

func (s *syncService) ProcessDocument(doc map[string]interface{}) error {
	docType, _ := doc["type"].(string)
	if docType != "occurrence" {
		return nil
	}

	jsonBytes, _ := json.Marshal(doc)
	
	var data IncomingOccurrenceData
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Classification
		if data.ClassificationData.ClassificationID != "" {
			// Map -> JSON String 変換
			classJSON, _ := json.Marshal(data.ClassificationData.ClassClassification)
			cls := entity.ClassificationJSON{
				ClassificationID:    data.ClassificationData.ClassificationID,
				ClassClassification: string(classJSON),
			}
			if err := tx.Save(&cls).Error; err != nil { return err }
		}

		// 2. Place
		if data.PlaceData.PlaceID != "" {
			// Map -> JSON String 変換
			coordJSON, _ := json.Marshal(data.PlaceData.Coordinates)
			
			// Pointer -> Value 変換
			accuracy := 0.0
			if data.PlaceData.Accuracy != nil {
				accuracy = *data.PlaceData.Accuracy
			}

			pl := entity.Place{
				PlaceID:     data.PlaceData.PlaceID,
				PlaceNameID: "", // nilではなく空文字
				Coordinates: string(coordJSON),
				Accuracy:    accuracy,
			}
			if err := tx.Save(&pl).Error; err != nil { return err }
		}

		// 3. Occurrence
		createdAt, _ := time.Parse(time.RFC3339, data.CreatedAt)

		// String -> Int64 変換
		wsID, _ := strconv.ParseInt(data.WorkstationID, 10, 64)
		userID, _ := strconv.ParseInt(data.CreatedByUserID, 10, 64)

		// Pointer -> Value 変換
		projectID := ""
		if data.ProjectID != nil {
			projectID = *data.ProjectID
		}
		
		bodyLength := 0.0
		if data.OccurrenceData.BodyLength != nil {
			bodyLength = *data.OccurrenceData.BodyLength
		}

		langID := ""
		if data.LanguageID != nil {
			langID = *data.LanguageID
		}

		occ := entity.Occurrence{
			OccurrenceID:     data.ID,
			WorkstationID:    wsID,
			UserID:           userID,
			ProjectID:        projectID,
			IndividualID:     data.OccurrenceData.IndividualID,
			Lifestage:        data.OccurrenceData.Lifestage,
			Sex:              data.OccurrenceData.Sex,
			BodyLength:       bodyLength,
			Note:             data.OccurrenceData.Note,
			ClassificationID: data.ClassificationData.ClassificationID,
			PlaceID:          data.PlaceData.PlaceID,
			LanguageID:       langID,
			CreatedAt:        createdAt,
			Timezone:         data.Timezone,
		}
		if err := tx.Save(&occ).Error; err != nil { return err }

		log.Printf("Synced occurrence: %s", occ.OccurrenceID)
		return nil
	})
}

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
