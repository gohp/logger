package main

import (
	logde "logdemo"
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	c := logde.New()
	c.SetEncoding("json")
	c.CloseConsoleDisplay()
	logger := c.InitLogger()
	for i := 0; i < b.N; i++ {
		logger.Info("info level test", logger.With("Trace", "123445"))
	}
}

func BenchmarkLoggerWithFile(b *testing.B) {
	c := logde.New()
	c.SetDivision("size")
	c.SetTimeUnit(logde.Minute)
	c.SetEncoding("json")
	c.CloseConsoleDisplay()
	c.SetInfoFile("./logs/server.log")
	//c.SetErrorFile("./logs/server_err.log")
	logger := c.InitLogger()
	for i := 0; i < b.N; i++ {
		logger.Info("info level test", logger.With("Trace", "123445"))
	}
}
