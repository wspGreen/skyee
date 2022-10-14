package slog

import (
	"fmt"
	"runtime"
	"strings"
)

type LogLevel int

const (
	LogLevelAllOn  LogLevel = iota //开放说有日志
	LogLevelDebug                  //调试信息
	LogLevelInfo                   //资讯讯息
	LogLevelWarn                   //警告状况发生
	LogLevelError                  //一般错误，可能导致功能不正常
	LogLevelFatal                  //严重错误，会导致进程退出
	LogLevelAllOff                 //关闭所有日志
)

type ILogger interface {
	Write(str string)
}

type ISLog interface {
	SetLogger(logger ILogger) bool
	Level() LogLevel
	SetLevel(level LogLevel)

	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

var log ISLog

func SetLog(l ISLog) {
	log = l
}

func Log() ISLog {
	return log
}

func SetLogger(logger ILogger) bool {
	return log.SetLogger(logger)
}

func Info(v ...interface{}) {
	log.Info(v...)
}

func Debug(v ...interface{}) {
	log.Debug(v...)
}

func Error(v ...interface{}) {
	log.Error(v...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Warn(v ...interface{}) {
	log.Warn(v...)
}

func Stack() {
	buf := make([]byte, 1<<12)
	Error(string(buf[:runtime.Stack(buf, false)]))
}

func SimpleStack() string {
	_, file, line, _ := runtime.Caller(2)
	i := strings.LastIndex(file, "/") + 1
	i = strings.LastIndex((string)(([]byte(file))[:i-1]), "/") + 1

	return fmt.Sprintf("%s:%d", (string)(([]byte(file))[i:]), line)
}
