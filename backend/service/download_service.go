package service

import (
	"JetGet/backend/config/db"
	"JetGet/backend/pget"
	"JetGet/backend/util"
	"context"
	"fmt"
	"gorm.io/gorm"
	"os"
)

type DownloadService struct {
	ctx    context.Context
	db     *gorm.DB
	hostId db.HostId
}

func NewDownloadService(db *gorm.DB, hostId db.HostId) *DownloadService {
	return &DownloadService{db: db, hostId: hostId}
}

func (d *DownloadService) Startup(ctx context.Context) {
	d.ctx = ctx
}

var Version string

// DownloadFile 下载文件
func (d *DownloadService) DownloadFile(url string) (string, error) {
	// 这里实现文件下载逻辑
	args := util.ToPgetArgs(url)
	go func() {
		cli := pget.New()
		if err := cli.Run(context.Background(), Version, args); err != nil {
			if cli.Trace {
				fmt.Fprintf(os.Stderr, "Error:\n%+v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Error:\n  %v\n", err)
			}
		}
	}()
	return fmt.Sprintf("Downloading: %s", url), nil
}
