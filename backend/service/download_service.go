package service

import (
	"JetGet/backend/pget"
	"JetGet/backend/util"
	"context"
	"fmt"
	"os"
)

type DownloadService struct {
	ctx context.Context
}

func NewDownloadService() *DownloadService {
	return &DownloadService{}
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
