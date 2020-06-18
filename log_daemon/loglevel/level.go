package loglevel

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"time"
)

type loggerlevel int64

// UNKNOW  设置日志级别常量
const (
	UNKNOW loggerlevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

// Logger 定义日志结构体
type Logger struct {
	level loggerlevel
}

// 将字符串转换成loggerlevel类型
func parseLogger(s string) (loggerlevel, error) {
	switch s {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKNOW, err
	}
}

// 将loggerlevel转换成字符串类型
func parseLogger2(s loggerlevel) string {
	switch s {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "DEBUG"
	}
}

// NewLogger  构造函数
func NewLogger(s string) Logger {
	level, err := parseLogger(s)
	if err != nil {
		panic(err)
	}
	return Logger{
		level: level,
	}
}

// enable 判断日志是否够级别
func (l *Logger) enable(x loggerlevel) bool {
	return l.level <= x
}

//获取调用日志的函数信息
func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	return
}

func (l *Logger) log(lv loggerlevel, format string) {
	if l.enable(lv) {
		now := time.Now()
		lvstr := parseLogger2(lv)
		funcName, fileName, lineNo := getInfo(3)
		fmt.Printf("[%v]  [%v] [%v %v %v] %v\n", now.Format("2006-01-02 03:04:05"), lvstr, fileName, lineNo, funcName, format)
	}
}

// Debug 日志级别
func (l *Logger) Debug(format string) {
	l.log(DEBUG, format)
}

// Info 日志级别
func (l *Logger) Info(format string) {
	l.log(INFO, format)
}

// Warning 日志级别
func (l *Logger) Warning(format string) {
	l.log(WARNING, format)
}

// Error 日志级别
func (l *Logger) Error(format string) {
	l.log(ERROR, format)
}
