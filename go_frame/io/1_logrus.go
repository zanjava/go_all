package io

import (
	"fmt"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func InitLogrus(logFile, level string) *logrus.Logger {
	logger := logrus.New()
	switch strings.ToLower(level) {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		panic(fmt.Errorf("invalid log level %s", level))
	}

	//普通文本格式
	// logger.SetFormatter(&logrus.TextFormatter{
	// 	// ForceColors: true,  //强制显示颜色（仅在某些终端能正常工作）

	// 	DisableColors:   true,                      //强制不显示颜色
	// 	TimestampFormat: "2006-01-02 15:04:05.000", // 显示ms
	// })
	//json格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000", // 显示ms
	})

	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H",                      //指定日志文件的路径和名称，路径不存在时会创建
		rotatelogs.WithLinkName(logFile),         //为最新的一份日志创建软链接
		rotatelogs.WithRotationTime(1*time.Hour), //每隔1小时生成一份新的日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),    //只留最近7天的日志，或使用WithRotationCount只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(fout) //设置日志文件
	// logger.SetOutput(os.Stdout)  //日志输出到终端
	logger.SetReportCaller(true) //输出是从哪里调起的日志打印，日志里会包含func和file

	logger.AddHook(&AppHook{AppName: "zgw"}) //在输出日志【之前】执行钩子

	return logger
}

// 实现logrus.Hook接口
type AppHook struct {
	AppName string
}

// 适用于哪些Level
func (h *AppHook) Levels() []logrus.Level {
	// return logrus.AllLevels
	return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}

// 在Fire函数时可读取或修改logrus.Entry
func (h *AppHook) Fire(entry *logrus.Entry) error {
	entry.Data["app"] = h.AppName //修改logrus.Entry
	fmt.Println(entry.Message)    //读取logrus.Entry。比如将 Error、Fatal 和 Panic 级别的错误日志发送到 logstash、kafka 等
	return nil
}
