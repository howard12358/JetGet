package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

// SLog 全局的 SugaredLogger，方便在应用的任何地方直接使用
var SLog *zap.SugaredLogger

var CustomLogger *CustomWailsLogger

// CustomWailsLogger 结构体实现了 wails.Logger 接口
type CustomWailsLogger struct {
	logger *zap.SugaredLogger
}

// ---- 实现 wails.Logger 接口的所有方法 ----

func (l *CustomWailsLogger) Print(message string) {
	l.logger.Info(message)
}

// Zap 没有 Trace 级别, 我们将其映射到 Debug 级别
func (l *CustomWailsLogger) Trace(message string) {
	l.logger.Debug(message)
}

func (l *CustomWailsLogger) Debug(message string) {
	l.logger.Debug(message)
}

func (l *CustomWailsLogger) Info(message string) {
	l.logger.Info(message)
}

func (l *CustomWailsLogger) Warning(message string) {
	l.logger.Warn(message)
}

func (l *CustomWailsLogger) Error(message string) {
	l.logger.Error(message)
}

// Zap 的 Fatal 会直接导致程序退出 (os.Exit), Wails 的 Fatal 只是记录日志
// 为了避免程序直接退出，我们使用 DPanic。它在开发环境下会 panic，但在生产环境下只会记录 Error 日志。
func (l *CustomWailsLogger) Fatal(message string) {
	l.logger.DPanic(message)
}

// InitLogger 是日志系统的总入口和初始化函数
func InitLogger() *zap.SugaredLogger {
	// 1. 配置日志文件写入器 (使用 lumberjack 进行轮转)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   getLogFilePath(), // 日志文件路径
		MaxSize:    20,               // 每个日志文件最大尺寸 (MB)
		MaxBackups: 5,                // 保留的旧日志文件最大数量
		MaxAge:     30,               // 旧日志文件最长保留天数
		Compress:   true,             // 是否压缩旧日志文件
	})

	// 2. 配置控制台写入器
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 3. 配置日志编码器 (Encoder)
	// 开发环境使用易于阅读的 ConsoleEncoder
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 时间格式: "2006-01-02T15:04:05.000Z0700"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 在控制台输出带颜色的大写日志级别

	// 4. 创建 Core，它决定了日志的最低级别和输出位置
	// 我们创建一个Tee Core，它可以同时向多个目标写入日志
	core := zapcore.NewTee(
		// 写入文件的 Core：只记录 Info 及以上级别的日志
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), fileWriter, zap.InfoLevel),
		// 写入控制台的 Core：记录所有级别的日志 (Debug)
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, zap.DebugLevel),
	)

	// 5. 创建 Logger 实例
	// zap.AddCaller() 会在日志中添加调用日志函数的文件名和行号
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 6. 将 Logger 转换为 SugaredLogger 并赋值给全局变量
	SLog = logger.Sugar()
	CustomLogger = &CustomWailsLogger{logger: SLog}
	return SLog
}

// getLogFilePath 辅助函数，用于获取日志文件的标准存储位置
func getLogFilePath() string {
	// 建议将日志文件放在用户配置目录下，而不是程序安装目录
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		// 如果获取失败，则降级到当前目录
		return "app.log"
	}
	logDir := filepath.Join(userConfigDir, "JetGet") // 替换为你的应用名
	// 确保目录存在
	_ = os.MkdirAll(logDir, os.ModePerm)
	return filepath.Join(logDir, "app.log")
}
