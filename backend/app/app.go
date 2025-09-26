package app

import (
	"JetGet/backend/service"
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx             context.Context
	SysService      *service.SysService
	DownloadService *service.DownloadService
}

// NewApp creates a new App application struct
func NewApp(sysService *service.SysService, downloadService *service.DownloadService) *App {
	return &App{
		SysService:      sysService,
		DownloadService: downloadService,
	}
}

// Startup is called at application Startup
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.SysService.Startup(ctx)
	a.SysService.Startup(ctx)
}

// DomReady is called after front-end resources have been loaded
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
}

// BeforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// Shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
