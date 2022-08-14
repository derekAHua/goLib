package rmq

import (
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

type rmqLogger struct {
	initOnce sync.Once
	verbose  bool
}

func (r *rmqLogger) Level(level string) {}

func (r *rmqLogger) isVerbose() bool {
	r.initOnce.Do(func() {
		if os.Getenv("RMQ_SDK_VERBOSE") != "" {
			r.verbose = true
		} else {
			r.verbose = false
		}
	})
	return r.verbose
}

func (r *rmqLogger) getFields(fields map[string]interface{}) []zapcore.Field {
	var f = []zlog.Field{
		logSDKTopic,
	}
	for k, v := range fields {
		f = append(f, zap.Reflect(k, v))
	}
	return f
}

func (r *rmqLogger) Debug(msg string, fields map[string]interface{}) {
	if r.isVerbose() {
		zlog.DebugLogger(nil, zlog.LogNameRMQ, msg, r.getFields(fields)...)
	}
}

func (r *rmqLogger) Info(msg string, fields map[string]interface{}) {
	if r.isVerbose() {
		zlog.InfoLogger(nil, zlog.LogNameRMQ, msg, r.getFields(fields)...)
	}
}

func (r *rmqLogger) Warning(msg string, fields map[string]interface{}) {
	if r.isVerbose() {
		zlog.WarnLogger(nil, zlog.LogNameRMQ, msg, r.getFields(fields)...)
	}
}

func (r *rmqLogger) Error(msg string, fields map[string]interface{}) {
	zlog.ErrorLogger(nil, zlog.LogNameRMQ, msg, r.getFields(fields)...)
}

func (r *rmqLogger) Fatal(msg string, fields map[string]interface{}) {
	zlog.FatalLogger(nil, zlog.LogNameRMQ, msg, r.getFields(fields)...)
}

func (r *rmqLogger) OutputPath(path string) (err error) {
	return nil
}
