package log

import (
	"github.com/superOTAKU/onlineTexasPoker/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger //业务用

func InitLogger(logConfig *config.LogConfig) error {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{}
	config.ErrorOutputPaths = []string{}
	logFile := false
	if len(logConfig.FilePath) > 0 {
		config.OutputPaths = append(config.OutputPaths, logConfig.FilePath)
		config.ErrorOutputPaths = append(config.ErrorOutputPaths, logConfig.FilePath)
		logFile = true
	}
	if logConfig.Console || !logFile {
		config.OutputPaths = append(config.OutputPaths, "stdout")
		config.ErrorOutputPaths = append(config.ErrorOutputPaths, "stderr")
	}
	l, err := config.Build()
	if err != nil {
		return err
	}
	logger = l
	return nil
}

func Logger() *zap.Logger {
	if logger == nil { // 使用之前必须先初始化
		panic("logger is nil")
	}
	return logger
}
