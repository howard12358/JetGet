//go:build wireinject
// +build wireinject

package injector

import (
	"JetGet/backend/app"
	"JetGet/backend/config/db"
	"JetGet/backend/service"
	"github.com/google/wire"
)

// InitializeApp 使用 Wire 来构建和初始化整个应用。
// 它会返回一个包含所有服务和配置的 App 结构体。
func InitializeApp() (*app.App, error) {
	// wire.Build 会根据下面的 Provider 集合自动生成依赖注入代码
	wire.Build(
		// 数据库层 Provider
		db.WireDBSet,
		// 服务层 Provider
		service.ServicesSet,
		// 应用层 Provider
		app.NewApp,
	)
	// 返回值只是一个占位符，Wire 会在生成的文件中实现它
	return nil, nil
}
