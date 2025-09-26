package db

import (
	"JetGet/backend/storage/m"
	"JetGet/backend/util"
	"fmt"
	"github.com/google/wire"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"time"
)

// HostId 定义为一种类型，以便 Wire 可以区分不同的 string 依赖
type HostId string

// ProvideHostID 是一个 Provider，用于生成并提供 HostId
// 它现在是一个纯函数，没有副作用，不依赖任何外部状态。
func ProvideHostID() (HostId, error) {
	id, err := util.GenerateHostID()
	if err != nil {
		return "", fmt.Errorf("failed to generate host id: %w", err)
	}
	// 直接生成并返回，不存储任何状态
	return HostId(id), nil
}

// ProvideDB 是一个 Provider，用于创建并提供 GORM DB 实例
// 它依赖 HostId 来初始化数据库（尽管在此代码中没有直接使用，但保持依赖关系是好的实践）
func ProvideDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000", getDbPath())

	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	if err := gdb.AutoMigrate(&m.SysConfig{}); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

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

// WireDBSet 将数据库相关的 Provider 组合在一起
var WireDBSet = wire.NewSet(ProvideDB, ProvideHostID)
