package service

import (
	"JetGet/backend/config/db"
	"JetGet/backend/pget"
	"JetGet/backend/types/e"
	"JetGet/backend/types/m"
	"JetGet/backend/util"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type DownloadService struct {
	ctx    context.Context
	db     *gorm.DB
	hostId db.HostId

	sysService *SysService
}

func NewDownloadService(db *gorm.DB, hostId db.HostId, sysService *SysService) *DownloadService {
	return &DownloadService{db: db, hostId: hostId, sysService: sysService}
}

func (d *DownloadService) Startup(ctx context.Context) {
	d.ctx = ctx
}

var Version string

// DownloadFile 下载文件
func (d *DownloadService) DownloadFile(url string) (string, error) {
	taskId := uuid.New().String()
	resp, err := DoHeadRequest(url, d.sysService.GetProxy())
	if err != nil {
		return "", err
	}
	fileName := GetFileName(resp)
	savePath := d.sysService.GetDownloadPath()
	task := &m.DownloadTask{
		ID:       taskId,
		URL:      url,
		SavePath: savePath,
		FileName: fileName,
		Status:   m.StatusPending,
	}

	err = d.db.Create(task).Error
	if err != nil {
		log.Printf("Error creating task in DB: %v\n", err)
		return "", err
	}

	// 马上通知前端有一个新任务被创建了
	runtime.EventsEmit(d.ctx, e.DownloadNew, task)

	args := util.ToPgetArgs(url)
	go func() {
		cli := pget.New()
		cli.ProgressFn = func(downloaded, total, speed int64) {
			d.db.Model(&m.DownloadTask{}).Where("id = ?", taskId).Updates(m.DownloadTask{
				DownloadedSize: downloaded,
				TotalSize:      total,
				Status:         m.StatusDownloading,
			})

			runtime.EventsEmit(d.ctx, e.DownloadProgress, e.Progress{
				ID:         taskId,
				Downloaded: downloaded,
				Total:      total,
				Speed:      speed,
				Status:     m.StatusDownloading,
			})
		}
		if err := cli.Run(context.Background(), Version, args); err != nil {
			// 下载失败
			errMsg := fmt.Sprintf("Download failed: %v", err)
			log.Println(errMsg)
			d.db.Model(&m.DownloadTask{}).Where("id = ?", taskId).Updates(m.DownloadTask{
				Status:       m.StatusFailed,
				ErrorMessage: errMsg,
			})
			runtime.EventsEmit(d.ctx, e.DownloadFailed, e.Progress{
				ID:     taskId,
				Status: m.StatusFailed,
				Msg:    errMsg,
			})
		} else {
			// 下载成功
			log.Println("Download completed successfully, id:", taskId)
			d.db.Model(&m.DownloadTask{}).Where("id = ?", taskId).Updates(m.DownloadTask{
				Status:      m.StatusCompleted,
				CompletedAt: time.Now(),
			})
			// 通知前端任务完成
			runtime.EventsEmit(d.ctx, e.DownloadCompleted, e.Progress{
				ID:     taskId,
				Status: m.StatusCompleted,
			})
		}
	}()
	return fmt.Sprintf("taskId: %s", taskId), nil
}

func DoHeadRequest(url, proxy string) (*http.Response, error) {
	// 查询文件大小
	client := pget.NewClientByProxy(16, proxy)
	r, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Println("new request failed:", err)
	}
	res, err := client.Do(r)
	if err != nil {
		log.Println("failed to head request:", err)
		return res, err
	}

	if res.Header.Get("Accept-Ranges") != "bytes" {
		return res, errors.New("does not support range request")
	}
	if res.ContentLength <= 0 {
		return res, errors.New("invalid content length")
	}
	return res, nil
}

func GetFileName(resp *http.Response) string {
	// 优先检查 Content-Disposition header
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		// 使用 mime 包来正确解析可能包含特殊字符的 header
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			// filename* 优先，因为它支持 UTF-8
			if filename, ok := params["filename*"]; ok {
				// filename* 的格式是 charset''encoded-value
				if parts := strings.SplitN(filename, "''", 2); len(parts) == 2 {
					decodedFilename, err := url.QueryUnescape(parts[1])
					if err == nil {
						return decodedFilename
					}
				}
			}
			// 其次是普通的 filename
			if filename, ok := params["filename"]; ok {
				return filename
			}
		}
	}
	// 如果 header 不存在，则从最终的 URL 路径中提取 (使用 resp.Request.URL 可以获取重定向之后的最终 URL)
	finalURL := resp.Request.URL.Path
	filename := filepath.Base(finalURL)
	if filename != "." && filename != "/" {
		return filename
	}
	return "unrecognized"
}

// 用于 wails 生成前端实体
func (d *DownloadService) GenDownloadTaskResp(task *m.DownloadTaskResp) {}
func (d *DownloadService) GenProgress(progress *e.Progress)             {}
