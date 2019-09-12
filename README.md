# logger

基于Zap,可选日志文件归档方式

## TODO

- [x] 根据info/warn级别切割日志文件
- [x] 根据文件大小归档
- [x] 根据时间归档
- [x] 时间切割单元可选
- [x] Benchmark test

## Usage

- install logger with go get

`go get -u github.com/wzyonggege/logger`

1. 新建logger
```go
import (
 "github.com/wzyonggege/logger"
)

...

c := logger.New()
c.SetDivision("time")	    // 设置归档方式，"time"时间归档 "size" 文件大小归档，文件大小等可以在配置文件配置
c.SetTimeUnit(logger.Minute) // 时间归档 可以设置切割单位
c.SetEncoding("json")	    // 输出格式 "json" 或者 "console"

c.SetInfoFile("./logs/server.log")		// 设置info级别日志
c.SetErrorFile("./logs/server_err.log")	// 设置warn级别日志

c.InitLogger()
```

2. 从配置文件中加载(Toml,Yaml,Json)

```go
// toml file
c := logger.NewFromToml(confPath)

// yaml file
c := logger.NewFromYaml("configs/config.yaml")


// json file
c := logger.NewFromJson("configs/config.json")

c.InitLogger()
```

3. caller 

```go
c.SetCaller(true)
```

4. 输出

```go
import (
 "github.com/wzyonggege/logger"
)

...
c := logger.New()
c.InitLogger()

logger.Info("info level test")
logger.Error("error level test")
logger.Warn("warn level test")
logger.Debug("debug level test")
logger.Fatal("fatal level test")
```

```bash
{"level":"info","time":"2019-09-11T18:32:59.680+0800","msg":"info level test"}
{"level":"error","time":"2019-09-11T18:32:59.680+0800","msg":"error level test"}
{"level":"warn","time":"2019-09-11T18:32:59.681+0800","msg":"warn level test"}
{"level":"debug","time":"2019-09-11T18:32:59.681+0800","msg":"debug level test"}
{"level":"fatal","time":"2019-09-11T18:32:59.681+0800","msg":"fatal level test"}
```

5. with args
```go
import (
 "github.com/wzyonggege/logger"
)

...
c := logger.New()
c.InitLogger()

logger.Info("this is a log", logger.With("Trace", "12345677"))
logger.Info("this is a log", logger.WithError("error", errors.New("this is a new error")))
```

```bash
{"level":"info","time":"2019-09-11T18:38:51.022+0800","msg":"this is a log","Trace":"12345677"}
{"level":"info","time":"2019-09-11T18:38:51.026+0800","msg":"this is a log","error":"this is a new error"}
```

## Performance

```bash
BenchmarkLogger/logde_logger_without_fields-4            3000000               563 ns/op
BenchmarkLogger/logde_logger_with_fields-4               2000000               637 ns/op
BenchmarkLogger/logde_logger_without_fields_write_into_file-4             200000             13021 ns/op
BenchmarkLogger/logde_logger_with_fields_write_into_file-4                100000             12606 ns/op
```



