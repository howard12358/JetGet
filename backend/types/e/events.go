package e

import (
	"JetGet/backend/types/m"
)

const (
	DownloadNew       = "download:new"
	DownloadProgress  = "download:progress"
	DownloadCompleted = "download:completed"
	DownloadFailed    = "download:failed"
)

// ProgressEvent 定义了发送给前端的进度事件结构体
type ProgressEvent struct {
	ID         string       `json:"id"`
	Downloaded int64        `json:"downloaded"`
	Total      int64        `json:"total"`
	Speed      int64        `json:"speed"` // bytes per second
	Status     m.TaskStatus `json:"status"`
	Msg        string       `json:"msg"`
}
