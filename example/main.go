package main

import (
	"errors"
	"github.com/gohp/logger"
	"github.com/spf13/pflag"
	"time"
)

var (
	confPath string
)

func init() {
	pflag.StringVar(&confPath, "conf", "configs/config.toml", "default configs path")
}

func main() {
	c := logger.New()
	c.SetDivision("time")     // 设置归档方式，"time"时间归档 "size" 文件大小归档，文件大小等可以在配置文件配置
	c.SetTimeUnit(logger.Day) // 时间归档 可以设置切割单位
	c.SetEncoding("json")     // 输出格式 "json" 或者 "console"
	//c.Stacktrace = true

	c.SetInfoFile("./logs/server.log")      // 设置info级别日志
	c.SetErrorFile("./logs/server_err.log") // 设置warn级别日志
	c.InitLogger()

	logger.Info("info level test")
	logger.Error("dsdadadad level test", logger.WithError(errors.New("sabhksasas")))
	logger.Error("121212121212 error")
	logger.Warn("warn level test")
	logger.Debug("debug level test")

	time.Sleep(2 * time.Second) // 避免程序结束太快，没有上传sentry

	logger.Infof("info level test: %s", "111")
	logger.Errorf("error level test: %s", "111")
	logger.Warnf("warn level test: %s", "111")
	logger.Debugf("debug level test: %s", "111")

	logger.Info("this is a log", logger.With("Trace", "12345677"))
	logger.Info("this is a log", logger.WithError(errors.New("this is a new error")))
}
