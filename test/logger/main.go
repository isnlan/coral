package main

import (
	"context"

	"github.com/isnlan/coral/pkg/trace"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/isnlan/coral/pkg/logging"
)

func main() {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./foo.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
	}

	logging.Init("app", logging.NewWriteSyncerConfig(lumberJackLogger))

	logger := logging.MustGetLogger("mod1")
	logger.Info("only info")
	logger.With(trace.GetTraceFieldFrom(context.Background())...).Info("with http trace")

	logger.Info("of")
	Add()
}

func Add() {
	logger := logging.MustGetLogger("mod2")
	logger.Info("only info in mod2")
}
