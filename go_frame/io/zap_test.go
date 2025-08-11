package io_test

import (
	"go/frame/io"
	"testing"

	"go.uber.org/zap"
)

func TestZap1(t *testing.T) {
	logger := io.InitZap1("../log/zap.log")
	defer logger.Sync() //将缓存的内容同步到文件中
	logger.Debug("hello")
	logger.Info("hello", zap.Int("age", 18))
	logger.Error("hello", zap.Namespace("china"), zap.Int("age", 18))

	sugar := logger.Sugar() //通过Sugar进行格式化输出,但性能会低50%左右
	sugar.Infof("pi is %f", 3.14)
}

func TestZap2(t *testing.T) {
	logger := io.InitZap2("../log/zap2.log", "info")
	defer logger.Sync() //将缓存的内容同步到文件中
	logger.Debug("hello")
	logger.Info("hello", zap.Int("age", 18))
	logger.Error("hello", zap.Namespace("china"), zap.Int("age", 18))
}

// go test -v ./io -run=^TestZap1$ -count=1
// go test -v ./io -run=^TestZap2$ -count=1
