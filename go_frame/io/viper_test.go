package io_test

import (
	"fmt"
	"go/frame/io"
	"testing"
	"time"
)

func TestViper(t *testing.T) {
	//读取配置的第一种方式
	dbViper := io.InitViper("../conf", "mysql", io.YAML)
	dbViper.WatchConfig()          //确保在调用WatchConfig()之前添加了所有的配置路径(AddConfigPath)
	if dbViper.IsSet("blog.age") { //检查有没有此项配置
		age := dbViper.GetInt("blog.age")
		fmt.Println("age", age)
	} else {
		fmt.Println("blog.age不存在")
	}
	port := dbViper.GetInt("blog.port") //该项不存在时会返回0值
	fmt.Println("port", port)
	time.Sleep(20 * time.Second) //10秒之内修改一下配置文件，看看viper能不能读取最新值
	port = dbViper.GetInt("blog.port")
	fmt.Println("port", port)

	//读取配置的第二种方式
	logViper := io.InitViper("../conf", "log", io.YAML)
	type LogConfig struct {
		Level string `mapstructure:"level"` //Tag
		File  string `mapstructure:"file"`
	}
	var config LogConfig
	if err := logViper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(config.Level)
		fmt.Println(config.File)
	}
}

// go test -v ./io -run=^TestViper$ -count=1
