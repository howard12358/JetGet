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
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
)

type DownloadService struct {
	ctx    context.Context
	db     *gorm.DB
	hostId db.HostId
	log    *zap.SugaredLogger

	sysService *SysService
}

func NewDownloadService(db *gorm.DB, hostId db.HostId, sysService *SysService, log *zap.SugaredLogger) *DownloadService {
	return &DownloadService{db: db, hostId: hostId, sysService: sysService, log: log}
}

func (d *DownloadService) Startup(ctx context.Context) {
	d.ctx = ctx
}

var Version string

// DownloadFile 下载文件
func (d *DownloadService) DownloadFile(url string) (string, error) {
	taskId := uuid.New().String()
	resp, err := util.DoHeadRequest(url, d.sysService.GetProxy())
	if err != nil {
		return "", err
	}
	fileName := util.GetFileName(resp)
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
		d.log.Info("Error creating task in DB: %v\n", err)
		return "", err
	}

	// 通知前端有一个新任务被创建了
	runtime.EventsEmit(d.ctx, e.DownloadNew, task)

	args := util.ToPgetArgs(url)

	progPointer := &atomic.Pointer[e.Progress]{}
	cancelThrottle := d.setupDownloadRateThrottle(progPointer, taskId)

	go func() {
		defer cancelThrottle()

		cli := pget.New()
		cli.ProgressFn = func(downloaded, total, speed int64) {
			currProg := e.Progress{
				ID:         taskId,
				Downloaded: downloaded,
				Total:      total,
				Speed:      speed,
				Status:     m.StatusDownloading,
			}
			// 只有带速度的回调才更新 speed，否则保留上一次的速度
			if speed == -1 {
				if oldProg := progPointer.Load(); oldProg != nil {
					currProg.Speed = oldProg.Speed
				}
			}
			progPointer.Store(&currProg)
		}
		if err := cli.Run(context.Background(), Version, args); err != nil {
			// 下载失败
			errMsg := fmt.Sprintf("Download failed: %v", err)
			d.log.Errorln(errMsg)
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
			d.log.Infoln("Download completed successfully, id:", taskId)
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

// setupDownloadRateThrottle 初始化下载进度节流器
func (d *DownloadService) setupDownloadRateThrottle(
	progress *atomic.Pointer[e.Progress],
	taskId string,
) context.CancelFunc {
	// 创建一个 context 用于控制节流 goroutine 的生命周期
	throttleCtx, cancelThrottle := context.WithCancel(context.Background())

	// 启动一个独立的 goroutine 来处理节流后的更新
	go func() {
		// 每 200 毫秒触发一次
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-throttleCtx.Done():
				// 下载结束或失败时，这个 goroutine 会退出
				return
			case <-ticker.C:
				// 从 atomic.Value 中加载最新的进度
				prog := progress.Load()
				if prog == nil {
					continue
				}

				d.db.Model(&m.DownloadTask{}).Where("id = ?", taskId).Updates(m.DownloadTask{
					DownloadedSize: prog.Downloaded,
					TotalSize:      prog.Total,
					Status:         m.StatusDownloading,
				})
				runtime.EventsEmit(d.ctx, e.DownloadProgress, prog)
			}
		}
	}()
	return cancelThrottle
}

func (d *DownloadService) PageDownloadHistory(status string, page, size int) util.PaginatedResult {
	var (
		history []m.DownloadTaskResp
		total   int64
	)
	d.db.Model(&m.DownloadTask{}).Count(&total)

	tx := d.db.Scopes(util.Paginate(page, size))
	if status == string(m.StatusCompleted) {
		tx.Where("status = ?", m.StatusCompleted)
	} else {
		tx.Where("status <> ?", m.StatusCompleted)
	}
	tx.Find(&history)

	return util.PaginatedResult{
		List:  history,
		Pager: util.Pager{Page: page, PageSize: size, Total: total},
	}
}

// 用于 wails 生成前端实体
func (d *DownloadService) GenDownloadTaskResp(*m.DownloadTaskResp)  {}
func (d *DownloadService) GenProgress(*e.Progress)                  {}
func (d *DownloadService) GenPaginatedResult(*util.PaginatedResult) {}
func (d *DownloadService) GenPager(*util.Pager)                     {}
