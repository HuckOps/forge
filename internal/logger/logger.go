package logger

import (
	"bufio"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
}

func ScanOutput(logger *zap.Logger, reader io.Reader, stream string, fields ...zap.Field) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		switch stream {
		case "stdout":
			logger.Info(text, fields...)
		case "stderr":
			logger.Error(text, fields...)
		}
	}
}
