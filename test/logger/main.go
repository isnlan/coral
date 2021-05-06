package main

import (
	"errors"

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

	logging.Init(logging.NewFileConfig("blink", lumberJackLogger))

	logger := logging.MustGetLogger("mysvr")
	logger.Errorf("test err: %v", errors.New("my error"))

	Add()
}

func Add() {
	logger := logging.MustGetLogger("AddMod")
	logger.Info("ok")
}
