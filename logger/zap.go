package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	debugMode bool
	logDir    string
	logger    *zap.Logger
}

func convertInterfacesToZapFields(i []interface{}) []zap.Field {
	var ans []zap.Field
	for _, f := range i {
		ans = append(ans, f.(zap.Field))
	}
	return ans
}

func (l *zapLogger) Debug(msg string, fields ...interface{}) {
	fs := convertInterfacesToZapFields(fields)
	l.logger.Debug(msg, fs...)
}

func (l *zapLogger) Info(msg string, fields ...interface{}) {
	fs := convertInterfacesToZapFields(fields)
	l.logger.Info(msg, fs...)
}

func (l *zapLogger) Warn(msg string, fields ...interface{}) {
	fs := convertInterfacesToZapFields(fields)
	l.logger.Warn(msg, fs...)
}

func (l *zapLogger) Error(msg string, fields ...interface{}) {
	fs := convertInterfacesToZapFields(fields)
	l.logger.Error(msg, fs...)
}

func (l *zapLogger) Fatal(msg string, fields ...interface{}) {
	fs := convertInterfacesToZapFields(fields)
	l.logger.Fatal(msg, fs...)
}

// DefaultLogger initializing default logger
func NewZapLogger(debugMode bool, logDir string) (Logger, error) {
	logDir = strings.TrimSpace(logDir)

	var core zapcore.Core

	confg := zap.NewDevelopmentEncoderConfig()
	confg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//---------------------------------------------------------------------------
	// log enablers and conjunction
	//---------------------------------------------------------------------------
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	if debugMode {
		core = zapcore.NewTee(
			// stdout, stderr
			zapcore.NewCore(zapcore.NewJSONEncoder(confg), zapcore.Lock(zapcore.AddSync(os.Stderr)), highPriority),
			zapcore.NewCore(zapcore.NewJSONEncoder(confg), zapcore.Lock(zapcore.AddSync(os.Stdout)), lowPriority),
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(confg),
				zapcore.AddSync(colorable.NewColorableStdout()),
				zapcore.DebugLevel,
			),
		)
	} else {
		//---------------------------------------------------------------------------
		// errors logfile
		//---------------------------------------------------------------------------
		errFilepath := filepath.Join(logDir, "errors.log")
		errFile, err := os.OpenFile(errFilepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			return nil, fmt.Errorf("failed to create error log file %s: %s", errFilepath, err)
		}
		errFileLog := zapcore.Lock(zapcore.AddSync(errFile))

		//---------------------------------------------------------------------------
		// regular logfile
		//---------------------------------------------------------------------------
		stdFilepath := filepath.Join(logDir, "standard.log")
		stdFile, err := os.OpenFile(stdFilepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			return nil, fmt.Errorf("failed to create standard log file %s: %s", errFilepath, err)
		}
		stdFileLog := zapcore.Lock(zapcore.AddSync(stdFile))

		core = zapcore.NewTee(
			// files
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), errFileLog, highPriority),
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), stdFileLog, lowPriority),
		)
	}

	return &zapLogger{
		debugMode: debugMode,
		logDir:    logDir,
		logger:    zap.New(core),
	}, nil
}
