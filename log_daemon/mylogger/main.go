package main

import (
	"main/log_daemon/loglevel"
	"time"
)

//日志库
//需求分析
//1.支持往不同的地方输出日志, 上线开发的时候什么级别都输出，但上线之后只有INFO以上
//2.日志要分级别   debug  trace info warning error  Fatal
//3.日志要支持开关控制
//4.完整的日志记录要有时间，行号，文件名，日志级别，日志信息
//5.日志文件要切割

func main() {
	//log := loglevel.NewLogger("WARNING")
	log := loglevel.NewFileLogger("./", "dyx.log", "INFO")
	for {
		log.Debug("这是一条DEBUG日志")
		log.Info("这是一条INFO日志")
		log.Warning("这是一条WARNING日志")
		log.Error("这是一条ERROR日志")
		time.Sleep(time.Second * 3)
	}
}
