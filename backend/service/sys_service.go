package service

import (
	"JetGet/backend/config/db"
	"JetGet/backend/types/m"
	"JetGet/backend/util"
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysService struct {
	ctx    context.Context
	db     *gorm.DB
	hostId db.HostId
	log    *zap.SugaredLogger
}

// NewSysService 是 SysService 的构造函数，它接收 *gorm.DB 和 HostId 作为依赖
func NewSysService(db *gorm.DB, hostId db.HostId, log *zap.SugaredLogger) *SysService {
	return &SysService{
		db:     db,
		hostId: hostId,
		log:    log,
	}
}

func (s *SysService) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *SysService) ChooseDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(s.ctx, runtime.OpenDialogOptions{Title: "选择下载目录"})
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetConfig 获取系统配置
func (s *SysService) GetConfig() (*m.SysConfig, error) {
	var config m.SysConfig
	result := s.db.Where("id = ?", s.hostId).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newSysConfig := m.SysConfig{ID: string(s.hostId)}
			res := s.db.Model(&m.SysConfig{}).Create(&newSysConfig)
			if res.Error != nil {
				return nil, res.Error
			}
			return &newSysConfig, nil
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

	newConfig := config
	// 使用 FirstOrCreate 来简化逻辑：如果记录存在则加载，不存在则创建
	result := s.db.Where("id = ?", s.hostId).FirstOrCreate(&newConfig)
	if result.RowsAffected == 0 {
		result = s.db.Model(&m.SysConfig{}).Where("id = ?", s.hostId).Updates(config)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (s *SysService) GetDownloadPath() string {
	config, err := s.GetConfig()
	if err != nil {
		return util.DefaultDownloadDir()
	}
	return config.DownloadDir
}

func (s *SysService) GetProxy() string {
	config, err := s.GetConfig()
	if err != nil {
		return ""
	}
	return config.Proxy
}

// 用于 wails 生成前端实体
func (s *SysService) GenSysConfig(*m.SysConfig) {}
