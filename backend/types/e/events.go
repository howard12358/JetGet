package e

import (
	"JetGet/backend/types/m"
	"time"
)

const (
	DownloadNew       = "download:new"
	DownloadProgress  = "download:progress"
	DownloadCompleted = "download:completed"
	DownloadFailed    = "download:failed"
)

// Progress 定义了发送给前端的进度事件结构体
type Progress struct {
	ID          string       `json:"id"`
	Downloaded  int64        `json:"downloaded"`
	Total       int64        `json:"total"`
	Speed       int64        `json:"speed"` // bytes per second
	Status      m.TaskStatus `json:"status"`
	Msg         string       `json:"msg"`
	CompletedAt time.Time    `json:"completedAt"`
}
