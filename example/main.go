package main

import (
	"github.com/spf13/pflag"
	logde "logdemo"
)

var (
	confPath string
)

func init() {
	pflag.StringVar(&confPath, "conf", "config.toml", "default configs path")
}

func main() {
	// from config
	//c := logde.NewFromToml(confPath)
	c := logde.New()
	c.SetDivision("time")
	c.SetErrorFile("./logs/server_err.log")
	c.SetInfoFile("./logs/server.log")
	logger := c.InitLogger()

	logger.Info("info level test", map[string]string{"Trace": "1234567890"})
	logger.Error("error level test", map[string]string{"Trace": "1234567890"})
}
