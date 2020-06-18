package loglevel

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type fileLevel int64

// FileLogger 文件日志类型
type FileLogger struct {
	filename   string
	filepath   string
	fileObj    *os.File
	errFileObj *os.File
	level      loggerlevel
	logtime    string
}

// NewFileLogger  构造函数
func NewFileLogger(filePath, fileName, s string) *FileLogger {
	level, err := parseLogger(s)
	if err != nil {
		panic(err)
	}
	f1 := &FileLogger{
		filepath: filePath,
		filename: fileName,
		level:    level,
	}
	err = f1.initfile()
	if err != nil {
		panic(err)
	}
	return f1
}

// 创建，并打开日志文件
func (f *FileLogger) initfile() error {
	fullFileName := path.Join(f.filepath, f.filename)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open the logfile failed %v", err)
		return err
	}
	errFullFileName := path.Join(f.filepath, "err.log")
	errFileObj, err := os.OpenFile(errFullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open the errlogfile failed %v", err)
		return err
	}
	now := time.Now()
	f.logtime = now.Format("2006-01-02")
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}

// 检测日志级别
func (f *FileLogger) enable(x loggerlevel) bool {
	return f.level <= x
}

// 检查日期
func (f *FileLogger) checkDate() int {
	now := time.Now()
	date := now.Format("2020-01-02")
	check := strings.Compare(date, f.logtime)
	return check
}

//日志切割，按照日期对日志文件进行切割
func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	now := time.Now()
	nowstr := now.Format("2020-01-02")
	fileInfo, err := file.Stat()
	//2.将日志文件进行备份
	logName := path.Join(f.filepath, fileInfo.Name())      //当前日志文件的完整路径
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowstr) //备份日志文件的路径和名字
	os.Rename(logName, newLogName)
	//1.关闭当前的日志
	file.Close()
	//3.打开新的日志文件
	fileObj, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open  new log file failed, err:%v\n", err)
		return nil, err
	}
	//4.将打开的新日志文件对象赋值给f.fileobj
	return fileObj, nil
}

// 往日志文件输入数据
func (f *FileLogger) log(lv loggerlevel, format string) {
	if f.enable(lv) {
		now := time.Now()
		lvstr := parseLogger2(lv)
		funcName, fileName, lineNo := getInfo(3)
		if f.checkDate() != 0 {
			newfile, err := f.splitFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newfile
		}
		fmt.Fprintf(f.fileObj, "[%v]  [%v] [%v %v %v] %v\n", now.Format("2006-01-02 03:04:05"), lvstr, fileName, lineNo, funcName, format)

		if lv >= ERROR {
			if f.checkDate() != 0 {
				newfile, err := f.splitFile(f.errFileObj)
				if err != nil {
					return
				}
				f.errFileObj = newfile
			}
			fmt.Fprintf(f.errFileObj, "[%v]  [%v] [%v %v %v] %v\n", now.Format("2006-01-02 03:04:05"), lvstr, fileName, lineNo, funcName, format)

		}
	}
}

// Debug 日志级别
func (f *FileLogger) Debug(format string) {
	f.log(DEBUG, format)
}

// Info 日志级别
func (f *FileLogger) Info(format string) {
	f.log(INFO, format)
}

// Warning 日志级别
func (f *FileLogger) Warning(format string) {
	f.log(WARNING, format)
}

// Error 日志级别
func (f *FileLogger) Error(format string) {
	f.log(ERROR, format)
}
