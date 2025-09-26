package service

import "github.com/google/wire"

// ServicesSet 将所有服务的 Provider 组合在一起
var ServicesSet = wire.NewSet(
	NewSysService,
	NewDownloadService,
)
