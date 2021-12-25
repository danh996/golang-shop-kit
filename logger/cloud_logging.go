package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/logging"
)

type cloudLogger struct {
	client *logging.Client

	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

// NewCloudLogger constructor
func NewCloudLogger(ctx context.Context) (Logger, error) {
	client, err := logging.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("logging.NewClient: %w", err)
	}
	logID := os.Getenv("LOG_ID")
	return &cloudLogger{
		client,
		client.Logger(logID).StandardLogger(logging.Debug),
		client.Logger(logID).StandardLogger(logging.Info),
		client.Logger(logID).StandardLogger(logging.Warning),
		client.Logger(logID).StandardLogger(logging.Error),
		client.Logger(logID).StandardLogger(logging.Critical),
	}, nil
}

func (l *cloudLogger) makeTextPayload(msg string, fields ...interface{}) string {
	payload := msg
	for _, field := range fields {
		payload = fmt.Sprintf("%s\n%v", payload, field)
	}
	return payload
}

func (l *cloudLogger) Debug(msg string, fields ...interface{}) {
	l.debugLogger.Println(l.makeTextPayload(msg, fields...))
}

func (l *cloudLogger) Info(msg string, fields ...interface{}) {
	l.infoLogger.Println(l.makeTextPayload(msg, fields...))
}

func (l *cloudLogger) Warn(msg string, fields ...interface{}) {
	l.warnLogger.Println(l.makeTextPayload(msg, fields...))
}

func (l *cloudLogger) Error(msg string, fields ...interface{}) {
	l.errorLogger.Println(l.makeTextPayload(msg, fields...))
}

func (l *cloudLogger) Fatal(msg string, fields ...interface{}) {
	l.fatalLogger.Println(l.makeTextPayload(msg, fields...))
}
