package zlog

import (
	"go.uber.org/zap/zapcore"
)

func NewJsonEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &jsonEncoder{
		Encoder: zapcore.NewJSONEncoder(cfg),
	}
}

type jsonEncoder struct {
	zapcore.Encoder
}
