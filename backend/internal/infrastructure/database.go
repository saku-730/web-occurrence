package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq" // GORMは内部でこれを使うのだ
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB はデータベース接続を初期化し、*gorm.DBと*sql.DBインスタンスを返すのだ
func InitDB() (*gorm.DB, *sql.DB) {
	// DSN (Data Source Name) を環境変数から構築する
	// 例: "postgres://admin:password@localhost:5432/specimendb?sslmode=disable"
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// GORM用のロガー設定（開発中は詳しいログを出すのだ）
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // ログの出力先
		logger.Config{
			SlowThreshold: time.Second, // 遅いSQLの閾値
			LogLevel:      logger.Info, // 開発中は Info, 本番環境では Warn にするといいのだ
			Colorful:      true,        // ログをカラフルにする
		},
	)

	// GORMでデータベースに接続
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: newLogger, // 設定したロガーを適用
	})

	if err != nil {
		log.Fatalf("GORMデータベースへの接続に失敗しました: %v", err)
	}

	// GORMから標準ライブラリの *sql.DB インスタンスを取得するのだ
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("GORMから*sql.DBの取得に失敗しました: %v", err)
	}

	// 接続テスト
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("データベースへのPingに失敗しました: %v", err)
	}

	fmt.Println("データベースに正常に接続しましたのだ。")
	return db, sqlDB // GORMのDBと、標準のsql.DBの両方を返す
}
