package m

import "time"

type SysConfig struct {
	ID          string    `gorm:"column:id;type:varchar(32);primaryKey;not null"`
	DownloadDir string    `gorm:"column:download_dir;type:varchar(255);null"`
	Proxy       string    `gorm:"column:proxy;type:varchar(100);null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (sc SysConfig) TableName() string {
	return "sys_config"
}
