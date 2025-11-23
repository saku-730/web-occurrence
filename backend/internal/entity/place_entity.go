// saku-730/web-occurrence/web-occurrence-main/backend/internal/entity/place_entity.go

package entity

type Place struct {
	PlaceID     string  `json:"place_id" gorm:"primaryKey;column:place_id;type:text"`
	// ★修正: stringから*string（ポインタ）に変更し、NULLを許可するのだ
	PlaceNameID *string `json:"place_name_id" gorm:"column:place_name_id;type:text"` 
	Coordinates string  `json:"coordinates" gorm:"column:coordinates;type:jsonb"` // JSONBは一旦stringで受ける
	Accuracy    float64 `json:"accuracy" gorm:"column:accuracy"`
}

func (Place) TableName() string {
	return "places"
}
