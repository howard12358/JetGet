package service

import (
	"JetGet/backend/config/db"
	"JetGet/backend/storage/model"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SysService struct {
	ctx context.Context
}

func NewSysService() *SysService {
	return &SysService{}
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
func (s *SysService) GetConfig() (*model.SysConfig, error) {
	var config model.SysConfig
	result := db.DB.Where("id = ?", db.HostId).First(&config)
	if result.Error != nil {
		// 如果没有找到配置，创建一个新的
		config.ID = db.HostId
		config.DownloadDir = ""
		config.Proxy = ""
		return &config, nil
	}

	return &config, nil
}

// SaveConfig 保存系统配置
func (s *SysService) SaveConfig(downloadDir, proxy string) error {
	config := model.SysConfig{
		ID:          db.HostId,
		DownloadDir: downloadDir,
		Proxy:       proxy,
	}

	// 先尝试创建，如果失败则更新
	result := db.DB.Create(&config)
	if result.Error != nil {
		// 如果创建失败（可能是因为主键冲突），则更新现有记录
		result = db.DB.Model(&model.SysConfig{}).Where("id = ?", db.HostId).Updates(map[string]interface{}{
			"download_dir": downloadDir,
			"proxy":        proxy,
		})
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
