package service

import (
	"JetGet/backend/config/db"
	"JetGet/backend/storage/m"
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type SysService struct {
	ctx    context.Context
	db     *gorm.DB
	hostId db.HostId
}

// NewSysService 是 SysService 的构造函数，它接收 *gorm.DB 和 HostId 作为依赖
func NewSysService(db *gorm.DB, hostId db.HostId) *SysService {
	return &SysService{
		db:     db,
		hostId: hostId,
	}
}

func (s *SysService) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *SysService) ChooseDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "选择下载目录",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetConfig 获取系统配置
func (s *SysService) GetConfig() (*m.SysConfig, error) {
	var config m.SysConfig
	// 使用注入的 s.db 和 s.hostId
	result := s.db.Where("id = ?", s.hostId).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果没有找到配置，返回一个带默认值的新配置
			return &m.SysConfig{
				ID:          string(s.hostId),
				DownloadDir: "",
				Proxy:       "",
			}, nil
		}
		return nil, result.Error
	}
	return &config, nil
}

// SaveConfig 保存系统配置
func (s *SysService) SaveConfig(downloadDir, proxy string) error {
	config := m.SysConfig{
		ID:          string(s.hostId),
		DownloadDir: downloadDir,
		Proxy:       proxy,
	}
	// 使用 FirstOrCreate 来简化逻辑：如果记录存在则加载，不存在则创建
	result := s.db.Where("id = ?", s.hostId).FirstOrCreate(&m.SysConfig{})
	if result.RowsAffected == 0 {
		result = s.db.Model(&m.SysConfig{}).Where("id = ?", s.hostId).Updates(config)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
