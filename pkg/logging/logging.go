package logging

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/derror"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var accessLogger *zap.SugaredLogger
var ApplicationLogger *zap.Logger

func AccessLogging(request *http.Request, err error) {
	if err != nil {
		var appErr derror.ApplicationError
		if errors.As(err, &appErr) {
			switch appErr.Level {
			case "error":
				accessLogger.Errorw(appErr.Msg,
					zap.Int("statusCode", appErr.Code),
					zap.Error(appErr.Err),
					zap.String("errStack", fmt.Sprintf("%+v", err)),
					zap.String("host", request.Host),
					zap.String("remoteAddress", request.RemoteAddr),
					zap.String("method", request.Method),
					zap.String("path", request.URL.Path),
					zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())),
					zap.String("userID", dcontext.GetUserIDFromContext(request.Context())))
			case "warn":
				accessLogger.Warnw(appErr.Msg,
					zap.Int("statusCode", appErr.Code),
					zap.Error(appErr.Err),
					zap.String("errStack", fmt.Sprintf("%+v", err)),
					zap.String("host", request.Host),
					zap.String("remoteAddress", request.RemoteAddr),
					zap.String("method", request.Method),
					zap.String("path", request.URL.Path),
					zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())),
					zap.String("userID", dcontext.GetUserIDFromContext(request.Context())))
			}
		} else {
			accessLogger.Errorw(appErr.Msg,
				zap.Int("statusCode", appErr.Code),
				zap.Error(appErr.Err),
				zap.String("errStack", fmt.Sprintf("%+v", err)),
				zap.String("host", request.Host),
				zap.String("remoteAddress", request.RemoteAddr),
				zap.String("method", request.Method),
				zap.String("path", request.URL.Path),
				zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())),
				zap.String("userID", dcontext.GetUserIDFromContext(request.Context())))
		}
	} else {
		accessLogger.Infow("succeed in access",
			zap.Int("statusCode", http.StatusOK),
			zap.String("host", request.Host),
			zap.String("remoteAddress", request.RemoteAddr),
			zap.String("method", request.Method),
			zap.String("path", request.URL.Path),
			zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())),
			zap.String("userID", dcontext.GetUserIDFromContext(request.Context())))
	}
}

func NewAccessLogger(zapCoreLevel zapcore.Level, env string) {
	var config zap.Config
	if env == "production" {
		config = zap.Config{
			Encoding:         "json",
			Level:            zap.NewAtomicLevelAt(zapCoreLevel),
			OutputPaths:      []string{"pkg/logging/log/access.log"},
			ErrorOutputPaths: []string{"pkg/logging/log/access.log"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:  "message",
				LevelKey:    "level",
				EncodeLevel: zapcore.CapitalLevelEncoder,
				TimeKey:     "time",
				EncodeTime:  zapcore.ISO8601TimeEncoder,
			},
		}
	} else {
		config = zap.Config{
			Encoding: "json",
			Level:    zap.NewAtomicLevelAt(zapCoreLevel),
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:  "message",
				LevelKey:    "level",
				EncodeLevel: zapcore.CapitalLevelEncoder,
				TimeKey:     "time",
				EncodeTime:  zapcore.ISO8601TimeEncoder,
			},
		}
	}
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	accessLogger = logger.Sugar()
}

func ApplicationErrorLogging(request *http.Request, err error) {
	if err != nil {
		var appErr derror.ApplicationError
		if errors.As(err, &appErr) {
			switch appErr.Level {
			case "error":
				ApplicationLogger.Error(appErr.Msg,
					zap.Error(appErr.Err),
					zap.String("errStack", fmt.Sprintf("%+v", err)),
					zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())))
			case "warn":
				ApplicationLogger.Error(appErr.Msg,
					zap.Error(appErr.Err),
					zap.String("errStack", fmt.Sprintf("%+v", err)),
					zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())))
			}
		} else {
			ApplicationLogger.Error(appErr.Msg,
				zap.Error(appErr.Err),
				zap.String("errStack", fmt.Sprintf("%+v", err)),
				zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())))
		}
	}
}

func NewApplicationLogger(zapCoreLevel zapcore.Level, env string) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	var writeSyncer zapcore.WriteSyncer
	if env == "production" {
		file, err := os.OpenFile("pkg/logging/log/application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		writeSyncer = zapcore.AddSync(file)
	} else {
		file := os.Stdout
		writeSyncer = zapcore.AddSync(file)
	}

	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapCoreLevel,
	)

	return zap.New(zapcore.NewTee(
		consoleCore,
		logCore,
	))
}

func init() {
	env := os.Getenv("ENV")
	var zapCoreLevel zapcore.Level
	if env == "production" {
		zapCoreLevel = zap.InfoLevel
	} else {
		zapCoreLevel = zap.DebugLevel
	}
	NewAccessLogger(zapCoreLevel, env)
	ApplicationLogger = NewApplicationLogger(zapCoreLevel, env)
	ApplicationLogger = ApplicationLogger.WithOptions(zap.AddCaller())

	defer accessLogger.Sync()
	defer ApplicationLogger.Sync()
}
