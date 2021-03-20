package logging

import (
	"InfecShotAPI/pkg/dcontext"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var accessLogger *zap.SugaredLogger

func AccessLogging(request *http.Request, err error) {
	if err != nil {

	} else {
		accessLogger.Infow("incoming request",
			zap.String("host", request.Host),
			zap.String("remoteAddress", request.RemoteAddr),
			zap.String("method", request.Method),
			zap.String("path", request.URL.Path),
			zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())))
	}
}

func NewAccessLogger(zapCoreLevel zapcore.Level) {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapCoreLevel),
		OutputPaths:      []string{"pkg/logging/log/access.log"},
		ErrorOutputPaths: []string{"pkg/logging/log/access.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "time",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			StacktraceKey: "stackTrace",
		},
	}
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	accessLogger = logger.Sugar()
}

func init() {
	env := os.Getenv("ENV")
	var zapCoreLevel zapcore.Level
	if env == "production" {
		zapCoreLevel = zap.InfoLevel
	} else {
		zapCoreLevel = zap.DebugLevel
	}
	NewAccessLogger(zapCoreLevel)
	defer accessLogger.Sync()
}
