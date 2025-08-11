package io_test

import (
	"go/frame/io"
	"log/slog"
	"testing"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LOG      = "这里是日志内容"
	TIME_FMT = "2006-01-02 15:04:05.000"
)

func InitSlog() *slog.Logger {
	logFile := "../log/slog.log"
	//fout, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H",                      //指定日志文件的路径和名称，路径不存在时会创建
		rotatelogs.WithLinkName(logFile),         //为最新的一份日志创建软链接
		rotatelogs.WithRotationTime(1*time.Hour), //每隔1小时生成一份新的日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),    //只留最近7天的日志，或使用WithRotationCount只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}

	logger := slog.New(
		&io.SlogContextHandler{
			//json格式
			slog.NewJSONHandler(fout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key != slog.TimeKey {
						return a
					}
					t := a.Value.Time()
					a.Value = slog.StringValue(t.Format(TIME_FMT))

					return a
				},
			}),
		},
	)
	return logger
}

func InitLogrus() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	//JSON格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: TIME_FMT, // 显示ms
	})

	logFile := "../log/logrus.log"
	//fout, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H",                      //指定日志文件的路径和名称，路径不存在时会创建
		rotatelogs.WithLinkName(logFile),         //为最新的一份日志创建软链接
		rotatelogs.WithRotationTime(1*time.Hour), //每隔1小时生成一份新的日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),    //只留最近7天的日志，或使用WithRotationCount只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(fout)       //设置日志文件
	logger.SetReportCaller(true) //输出是从哪里调起的日志打印，日志里会包含func和file
	return logger
}

func InitZap() *zap.Logger {
	logFile := "../log/zap.log"
	//fout, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H",                      //指定日志文件的路径和名称，路径不存在时会创建
		rotatelogs.WithLinkName(logFile),         //为最新的一份日志创建软链接
		rotatelogs.WithRotationTime(1*time.Hour), //每隔1小时生成一份新的日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),    //只留最近7天的日志，或使用WithRotationCount只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(TIME_FMT) //指定时间格式
	encoderConfig.TimeKey = "time"                                   //默认是ts
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder          //指定level的显示样式
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), //日志为json格式
		zapcore.AddSync(fout),                 //指定输出到文件
		zapcore.InfoLevel,                     //设置最低级别
	)
	logger := zap.New(
		core,
		zap.AddCaller(), //上报文件名和行号
	)
	return logger
}

func BenchmarkLoggerSlog(b *testing.B) {
	logger := InitSlog()
	b.ResetTimer()
	for b.Loop() {
		logger.Info(LOG, "name", "zgw", "age", 18)
		// logger.Info(LOG, slog.String("name", "dqq"), slog.Int("age", 18))  //显式指定数据类型也没有带来性能的提升
	}
}

func BenchmarkLoggerLogrus(b *testing.B) {
	logger := InitLogrus()
	b.ResetTimer()
	for b.Loop() {
		logger.WithFields(logrus.Fields{"name": "zgw", "age": 18}).Info(LOG)
	}
}

func BenchmarkLoggerZap(b *testing.B) {
	logger := InitZap()
	b.ResetTimer()
	for b.Loop() {
		logger.Info(LOG, zap.String("name", "zgw"), zap.Int("age", 18))
	}
	logger.Sync() //把缓冲里的内容刷入磁盘
}

// go test ./io -bench=^BenchmarkLoggerSlog$ -run=^$
// go test ./io -bench=^BenchmarkLoggerLogrus$ -run=^$
// go test ./io -bench=^BenchmarkLoggerZap$ -run=^$
// go test ./io -bench=^BenchmarkLogger -run=^$

/**
goos: windows
goarch: amd64
pkg: go/frame/io
cpu: 11th Gen Intel(R) Core(TM) i5-1145G7 @ 2.60GHz
BenchmarkLoggerSlog-8             151940              7729 ns/op
BenchmarkLoggerLogrus-8            96739             13373 ns/op
BenchmarkLoggerZap-8              211273              6037 ns/op
*/
