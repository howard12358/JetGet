package util

import "gorm.io/gorm"

// Pager 分页参数
type Pager struct {
	Page     int   `json:"page" form:"page"`         // 页码
	PageSize int   `json:"pageSize" form:"pageSize"` // 每页大小
	Total    int64 `json:"total"`                    // 总行数
}

// PaginatedResult 分页返回结果
type PaginatedResult struct {
	List  interface{}    `json:"list"` // 数据列表
	Pager `json:"pager"` // 分页信息
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
