package core

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	InitLogger("log.txt")
}

// InitLogger 初始化 slog 日志，输出到指定文件（追加模式）。
func InitLogger(filename string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger = slog.Default()
		logger.Error("failed to open log file", "error", err)
		return
	}
	logger = slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// LogDebug 输出 Debug 级别日志。
func LogDebug(msg string, args ...any) {
	if logger != nil {
		logger.Debug(msg, args...)
	}
}
