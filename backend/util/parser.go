package util

import (
	"os"
	"path/filepath"
	"runtime"
)

func ToPgetArgs(url string) []string {
	var ags []string
	ags = append(ags, "-x")
	ags = append(ags, "http://127.0.0.1:7897")

	ags = append(ags, "-p")
	ags = append(ags, "4")
	ags = append(ags, "-o")
	ags = append(ags, defaultDownloadsDir())

	ags = append(ags, url)
	return ags
}

func defaultDownloadsDir() string {
	// 简单且通常有效的做法：用 home + "Downloads"
	// 更严格的实现可以在 Linux 读取 ~/.config/user-dirs.dirs 中 XDG_DOWNLOAD_DIR
	home, err := os.UserHomeDir()
	if err != nil {
		return "." // 兜底
	}

	switch runtime.GOOS {
	case "windows":
		// Windows 的 Downloads 通常在 %USERPROFILE%\Downloads
		return filepath.Join(home, "Downloads")
	case "darwin":
		return filepath.Join(home, "Downloads")
	default: // linux / other
		// 尝试读取 XDG 配置（更完善），这里先用 home/Downloads
		return filepath.Join(home, "Downloads")
	}
}
