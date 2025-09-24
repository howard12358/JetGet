package db

import (
	"JetGet/backend/storage/model"
	"JetGet/backend/util"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"time"
)

var DB *gorm.DB

var HostId string

func InitDB() error {
	hostId, err := util.GenerateHostID()
	if err != nil {
		log.Fatalf("failed to generate host id: %v", err)
		return err
	}
	HostId = hostId
	db, err := initSqlite()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
		return err
	}
	// 自动迁移
	if err := db.AutoMigrate(&model.SysConfig{}); err != nil {
		log.Fatalf("failed to migrate db: %v", err)
		return err
	}
	return nil
}

func initSqlite() (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000", getDbPath())

	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = gdb
	return gdb, nil
}

func getDbPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("failed to get user config dir: %v", err)
	}
	appDir := filepath.Join(dir, "JetGet")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		log.Fatalf("failed to create app dir: %v", err)
	}
	return filepath.Join(appDir, "jet_get.db")
}
