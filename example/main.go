package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gohp/logger"
	"github.com/spf13/pflag"
)

var (
	confPath string
)

func init() {
	pflag.StringVar(&confPath, "conf", "configs/config.toml", "default configs path")
}

func main() {
	c := logger.New()
	c.SetDivision("time")                   // 设置归档方式，"time"时间归档 "size" 文件大小归档，文件大小等可以在配置文件配置
	c.SetTimeUnit(logger.Day)               // 时间归档 可以设置切割单位
	c.SetEncoding("json")                   // 输出格式 "json" 或者 "console"
	c.Stacktrace = true                     // 添加 Stacktrace, 默认false
	c.SetInfoFile("./logs/server.log")      // 设置info级别日志
	c.SetErrorFile("./logs/server_err.log") // 设置warn级别日志
	c.SetEncodeTime("2006-01-02 15:04:05")  // 设置时间格式
	c.InitLogger()                          // 初始化

	logger.Info("info level test")
	logger.Error("dsdadadad level test", logger.WithError(errors.New("sabhksasas")))
	logger.Error("121212121212 error")
	logger.Warn("warn level test")
	logger.Debug("debug level test")
	logger.Infof("info level test: %s", "111")
	logger.Errorf("error level test: %s", "111")
	logger.Warnf("warn level test: %s", "111")
	logger.Debugf("debug level test: %s", "111")
	logger.Info("this is a log", logger.With("Trace", "12345677"))
	logger.Info("this is a log", logger.WithError(errors.New("this is a new error")))

	// AddContext for context.Context
	ctx := context.Background()
	ctx = logger.AddContext(ctx, logger.With("trace", "12345678"))
	logger.Ctx(ctx).Info("info logger with ctx value")
	logger.Ctx(ctx).Error("warn logger with ctx value")

	// GAddContext for gin
	g := gin.New()
	g.Use(func(c *gin.Context) {
		logger.GAddContext(c, logger.With("trace", "12345678"))
		c.Next()
	})
	g.GET("/", func(c *gin.Context) {
		logger.Ctx(c).Infof("get request url: %v", c.Request.URL)
		return
	})
	g.Run()
	// {"level":"info","time":"2022-09-27 16:23:08","msg":"get request: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"}
}
