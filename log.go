package main

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// 初始化日志模块
func initLogger() error {
	logFilePath := Config.Log.Path
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	// 解析日志记录的等级信息
	level, err := logrus.ParseLevel(Config.Log.Level)
	if err != nil {
		return err
	}
	// 初始化日志结构
	Log = &logrus.Logger{
		Out:       file,
		Level:     level,
		Formatter: new(logrus.JSONFormatter),
	}
	Log.Infof("InitLogger: path: %s, level: %s, formatter: json", Config.Log.Path, level)
	return nil
}
