package main

import (
	"JetGet/backend/injector"
	"JetGet/backend/zaplog"
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	app, err := injector.InitializeApp()
	if err != nil {
		zaplog.SLog.Fatalf("failed to init app: %v", err)
		return
	}
	// 确保在程序退出时，所有缓冲的日志都被写入
	defer zaplog.SLog.Sync()

	// Create application with options
	err = wails.Run(&options.App{
		Title:             "JetGet",
		Width:             1050,
		Height:            700,
		MinWidth:          1050,
		MinHeight:         700,
		MaxWidth:          1050,
		MaxHeight:         700,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer:       &assetserver.Options{Assets: assets},
		Menu:              nil,
		Logger:            zaplog.CustomLogger,
		LogLevel:          logger.DEBUG,
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
		},
		OnDomReady:       app.DomReady,
		OnBeforeClose:    app.BeforeClose,
		OnShutdown:       app.Shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app.SysService,
			app.DownloadService,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarDefault(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "JetGet",
				Message: "",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		zaplog.SLog.Fatalf("failed to run app: %v", err)
	}
}
