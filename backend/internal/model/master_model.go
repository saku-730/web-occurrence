package model

type Language struct {
	LanguageID     int64  `json:"language_id" gorm:"primaryKey;column:language_id"`
	LanguageShort  string `json:"language_short" gorm:"column:language_short"`
	LanguageCommon string `json:"language_common" gorm:"column:language_common"`
}

func (Language) TableName() string { return "languages" }

type FileType struct {
	FileTypeID int64  `json:"file_type_id" gorm:"primaryKey;column:file_type_id"`
	TypeName   string `json:"type_name" gorm:"column:type_name"`
}

func (FileType) TableName() string { return "file_types" }

type FileExtension struct {
	ExtensionID   int64  `json:"extension_id" gorm:"primaryKey;column:extension_id"`
	ExtensionText string `json:"extension_text" gorm:"column:extension_text"`
	FileTypeID    int64  `json:"file_type_id" gorm:"column:file_type_id"`
}

func (FileExtension) TableName() string { return "file_extensions" }

type UserRole struct {
	RoleID   int64  `json:"role_id" gorm:"primaryKey;column:role_id"`
	RoleName string `json:"role_name" gorm:"column:role_name"`
}

func (UserRole) TableName() string { return "user_roles" }

// WorkstationUser は users テーブルから必要な情報だけ抜粋したもの
type WorkstationUser struct {
	UserID      int64  `json:"user_id" gorm:"column:user_id"`
	DisplayName string `json:"display_name" gorm:"column:display_name"`
}
