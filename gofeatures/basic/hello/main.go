package main

import (
	"log/slog"
	"os"
	"own/simple/greetings"
)

const inPod = false

func main() {
	// 设置合适的logger
	slog.SetDefault(initLogger())

	// predefined logger
	msg, err := greetings.Hello("")
	if err != nil {
		slog.Error("failed", err)
	} else {
		slog.Info(msg)
	}
}

// initLogger 标准化slog输出, pod中使用json
func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: false, // 是否打印具体文件&输出的line
		Level:     slog.LevelDebug,
	}
	var handler slog.Handler
	if inPod {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}
