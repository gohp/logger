package logde

import (
	"github.com/BurntSushi/toml"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

const (
	TimeDivision = "time"
	SizeDivision = "size"
)

var Logger *Log

type Log struct {
	L *zap.Logger
}

type LogOptions struct {
	InfoFilename  string
	ErrorFilename string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	Division      string
	LevelSeparate bool
	TimeUnit      time.Duration
}

func infoLevel() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})
}

func warnLevel() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
}

func New() *LogOptions {
	return &LogOptions{
		Division:      "size",
		LevelSeparate: false,
	}
}

func NewFromToml(confPath string) *LogOptions {
	var c *LogOptions
	if _, err := toml.DecodeFile(confPath, &c); err != nil {
		panic(err)
	}
	return c
}

func (c *LogOptions) SetDivision(division string) {
	c.Division = division
}

func (c *LogOptions) SetErrorFile(path string) {
	c.LevelSeparate = true
	c.ErrorFilename = path
}

func (c *LogOptions) SetInfoFile(path string) {
	c.InfoFilename = path
}

// isOutput whether set output file
func (c *LogOptions) isOutput() bool {
	return c.InfoFilename != ""
}

func (c *LogOptions) InitLogger() *Log {
	var (
		core               zapcore.Core
		infoHook, warnHook io.Writer
		wsInfo             []zapcore.WriteSyncer
		wsWarn             []zapcore.WriteSyncer
	)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	wsInfo = append(wsInfo, zapcore.AddSync(os.Stdout))
	wsWarn = append(wsWarn, zapcore.AddSync(os.Stdout))

	// zapcore WriteSyncer setting
	if c.isOutput() {
		switch c.Division {
		case TimeDivision:
			infoHook = c.timeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.timeDivisionWriter(c.ErrorFilename)
			}
		case SizeDivision:
			infoHook = c.sizeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.sizeDivisionWriter(c.ErrorFilename)
			}
		}
		wsInfo = append(wsInfo, zapcore.AddSync(infoHook))
	}

	if c.ErrorFilename != "" {
		wsWarn = append(wsWarn, zapcore.AddSync(warnHook))
	}

	// Separate info and warning log
	if c.LevelSeparate {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
				zapcore.NewMultiWriteSyncer(wsInfo...),
				infoLevel()),
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
				zapcore.NewMultiWriteSyncer(wsWarn...),
				warnLevel()),
		)
	} else {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(wsInfo...),
			zap.InfoLevel)
	}

	// add caller
	caller := zap.AddCaller()
	// file line number display
	development := zap.Development()
	// init default key
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))

	logger := zap.New(core, caller, development)
	Logger = &Log{logger}
	return Logger
}

func (c *LogOptions) sizeDivisionWriter(filename string) io.Writer {
	hook := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxSize,
		Compress:   c.Compress,
	}
	return hook
}

func (c *LogOptions) timeDivisionWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithMaxAge(time.Duration(int64(24*time.Hour)*int64(c.MaxAge))),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func (logger *Log) Info(msg string, args map[string]string) {
	var fields []zap.Field
	for k, v := range args {
		fields = append(fields, zap.String(k, v))
	}
	logger.L.Info(msg, fields...)
}

func (logger *Log) Error(msg string, args map[string]string) {
	var fields []zap.Field
	for k, v := range args {
		fields = append(fields, zap.String(k, v))
	}
	logger.L.Error(msg, fields...)
}
