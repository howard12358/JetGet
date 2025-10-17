package m

import (
	"time"
)

// TaskStatus 定义任务状态的常量，便于管理 (这个保持不变)
type TaskStatus string

const (
	StatusPending     TaskStatus = "pending"
	StatusDownloading TaskStatus = "downloading"
	StatusPaused      TaskStatus = "paused"
	StatusCompleted   TaskStatus = "completed"
	StatusFailed      TaskStatus = "failed"
)

// DownloadTask 下载任务
type DownloadTask struct {
	ID             string     `gorm:"primaryKey;type:varchar(36)" json:"id"`                         // 任务的唯一ID, 使用 UUID
	URL            string     `gorm:"type:text;not null" json:"url"`                                 // 原始下载链接
	FileName       string     `gorm:"type:varchar(255);not null" json:"fileName"`                    // 文件名
	SavePath       string     `gorm:"type:text;not null" json:"savePath"`                            // 本地保存路径
	TotalSize      int64      `gorm:"type:bigint;default:0" json:"totalSize"`                        // 文件总大小 (Bytes)
	DownloadedSize int64      `gorm:"type:bigint;default:0" json:"downloadedSize"`                   // 已下载大小 (Bytes)
	Status         TaskStatus `gorm:"type:varchar(20);not null;default:pending;index" json:"status"` // 任务状态
	ErrorMessage   string     `gorm:"type:text" json:"errorMessage"`                                 // 失败时的错误信息
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"createdAt"`                               // 创建时间
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`                               // 最后更新时间
	CompletedAt    time.Time  `gorm:"index" json:"completedAt"`                                      // 完成时间
}

func (dt DownloadTask) TableName() string {
	return "download_task"
}

type DownloadTaskResp struct {
	Speed int64 `json:"speed"`
	DownloadTask
}
