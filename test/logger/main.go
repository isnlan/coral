package main

import (
	"errors"

	"go.uber.org/zap/zapcore"
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
	writeSyncer := zapcore.AddSync(lumberJackLogger)
	_ = `{"@timestamp": "%{time:2006-01-02 15:04:05.000 MST}", "service": "blink", "module": "[%{module}]", "func": "%{shortfunc}", "level": "%{level:.4s}", "msg": "%{message}"}`
	c := logging.Config{
		Format:  "json",
		LogSpec: "",
		Writer:  writeSyncer,
	}
	logging.Init(c)

	logger := logging.MustGetLogger("mysvr")
	logger.Errorf("test err: %v", errors.New("my error"))

	Add()
}

func Add() {
	logger := logging.MustGetLogger("AddMod")
	logger.Info("ok")
}
