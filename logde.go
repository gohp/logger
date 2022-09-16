package logger

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const (
	TimeDivision = "time"
	SizeDivision = "size"

	_defaultEncoding = "console"
	_defaultDivision = "size"
	_defaultUnit     = Hour
)

var (
	Logger                    *Log
	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		"console": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewConsoleEncoder(encoderConfig)
		},
		"json": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewJSONEncoder(encoderConfig)
		},
	}
)

type Log struct {
	L *zap.Logger
}

type LogOptions struct {
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding      string   `json:"encoding" yaml:"encoding" toml:"encoding"`
	InfoFilename  string   `json:"info_filename" yaml:"info_filename" toml:"info_filename"`
	ErrorFilename string   `json:"error_filename" yaml:"error_filename" toml:"error_filename"`
	MaxSize       int      `json:"max_size" yaml:"max_size" toml:"max_size"`
	MaxBackups    int      `json:"max_backups" yaml:"max_backups" toml:"max_backups"`
	MaxAge        int      `json:"max_age" yaml:"max_age" toml:"max_age"`
	Compress      bool     `json:"compress" yaml:"compress" toml:"compress"`
	Division      string   `json:"division" yaml:"division" toml:"division"`
	LevelSeparate bool     `json:"level_separate" yaml:"level_separate" toml:"level_separate"`
	TimeUnit      TimeUnit `json:"time_unit" yaml:"time_unit" toml:"time_unit"`
	Stacktrace    bool     `json:"stacktrace" yaml:"stacktrace" toml:"stacktrace"`
	//SentryConfig  SentryLoggerConfig `json:"sentry_config" yaml:"sentry_config" toml:"sentry_config"`
	closeDisplay int
	caller       bool
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
		Division:      _defaultDivision,
		LevelSeparate: false,
		TimeUnit:      _defaultUnit,
		Encoding:      _defaultEncoding,
		caller:        false,
	}
}

func NewFromToml(confPath string) *LogOptions {
	var c *LogOptions
	if _, err := toml.DecodeFile(confPath, &c); err != nil {
		panic(err)
	}
	return c
}

func NewFromYaml(confPath string) *LogOptions {
	var c *LogOptions
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return c
}

func NewFromJson(confPath string) *LogOptions {
	var c *LogOptions
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return c
}

func (c *LogOptions) SetDivision(division string) {
	c.Division = division
}

func (c *LogOptions) CloseConsoleDisplay() {
	c.closeDisplay = 1
}

func (c *LogOptions) SetCaller(b bool) {
	c.caller = b
}

func (c *LogOptions) SetTimeUnit(t TimeUnit) {
	c.TimeUnit = t
}

func (c *LogOptions) SetErrorFile(path string) {
	c.LevelSeparate = true
	c.ErrorFilename = path
}

func (c *LogOptions) SetInfoFile(path string) {
	c.InfoFilename = path
}

func (c *LogOptions) SetEncoding(encoding string) {
	c.Encoding = encoding
}

// isOutput whether set output file
func (c *LogOptions) isOutput() bool {
	return c.InfoFilename != ""
}

func (c *LogOptions) InitLogger() *Log {
	var (
		logger             *zap.Logger
		infoHook, warnHook io.Writer
		wsInfo             []zapcore.WriteSyncer
		wsWarn             []zapcore.WriteSyncer
	)

	if c.Encoding == "" {
		c.Encoding = _defaultEncoding
	}
	encoder := _encoderNameToConstructor[c.Encoding]

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

	if c.closeDisplay == 0 {
		wsInfo = append(wsInfo, zapcore.AddSync(os.Stdout))
		wsWarn = append(wsWarn, zapcore.AddSync(os.Stdout))
	}

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

	opts := make([]zap.Option, 0)
	cos := make([]zapcore.Core, 0)

	if c.LevelSeparate {
		cos = append(
			cos,
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), infoLevel()),
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsWarn...), warnLevel()),
		)
	} else {
		cos = append(
			cos,
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), zap.InfoLevel),
		)
	}

	opts = append(opts, zap.Development())

	if c.Stacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.WarnLevel))
	}

	if c.caller {
		opts = append(opts, zap.AddCaller())
	}

	logger = zap.New(zapcore.NewTee(cos...), opts...)

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
		filename+c.TimeUnit.Format(),
		rotatelogs.WithMaxAge(time.Duration(int64(24*time.Hour)*int64(c.MaxAge))),
		rotatelogs.WithRotationTime(c.TimeUnit.RotationGap()),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Info(msg string, args ...zap.Field) {
	Logger.L.Info(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	Logger.L.Error(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	Logger.L.Warn(msg, args...)
}

func Debug(msg string, args ...zap.Field) {
	Logger.L.Debug(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	Logger.L.Fatal(msg, args...)
}

func Infof(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Info(logMsg)
}

func Errorf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Error(logMsg)
}

func Warnf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Warn(logMsg)
}

func Debugf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Debug(logMsg)
}

func Fatalf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Fatal(logMsg)
}

func With(k string, v interface{}) zap.Field {
	return zap.Any(k, v)
}

func WithError(err error) zap.Field {
	return zap.NamedError("error", err)
}

// 给指定的context添加字段（关键方法）
func AddContext(ctx *gin.Context, fields ...zap.Field) {
	l, _ := ctx.Get(strconv.Itoa(0))
	logArgs, ok := l.([]zap.Field)
	if ok || logArgs == nil {
		logArgs = append(logArgs, fields...)
		ctx.Set(strconv.Itoa(0), logArgs)
	}
}

// 从指定的context返回一个zap实例（关键方法）
func WithContext(ctx *gin.Context) *Log {
	if ctx == nil {
		return nil
	}
	l, _ := ctx.Get(strconv.Itoa(0))
	logArgs, _ := l.([]zap.Field)

	ctxLogger := &Log{
		L: Logger.L,
	}

	if len(logArgs) > 0 {
		ctxLogger.L = ctxLogger.L.With(logArgs...)
	}
	return ctxLogger
}

func (l *Log) Info(msg string, args ...zap.Field) {
	l.L.Info(msg, args...)
}

func (l *Log) Error(msg string, args ...zap.Field) {
	l.L.Error(msg, args...)
}

func (l *Log) Warn(msg string, args ...zap.Field) {
	l.L.Warn(msg, args...)
}

func (l *Log) Debug(msg string, args ...zap.Field) {
	l.L.Debug(msg, args...)
}

func (l *Log) Fatal(msg string, args ...zap.Field) {
	l.L.Fatal(msg, args...)
}

func (l *Log) Infof(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Info(logMsg)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Error(logMsg)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Warn(logMsg)
}

func (l *Log) Debugf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Debug(logMsg)
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Fatal(logMsg)
}
