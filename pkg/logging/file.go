package logging

import (
	"io"
	"time"

	"go.uber.org/zap/zapcore"
)

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func EpochTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.UTC().Format("2006-01-02T15:04:05Z07:00.000Z"))
}

func NewFileConfig(w io.Writer) Config {
	writeSyncer := zapcore.AddSync(w)
	c := Config{
		Format:  "json",
		LogSpec: "",
		Writer:  writeSyncer,
	}
	return c
}

type Rotator interface {
	Rotate() error
}
