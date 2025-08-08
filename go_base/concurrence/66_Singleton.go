package concurrence

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// 解析配置文件的单例模式
type config struct {
	Password      string
	ServerAddress string
}

var (
	cfg *config
	mux sync.Mutex
)

// GetConfig 获取配置文件的单例实例
// 双重检查锁定（Double-Checked Locking）模式
// 不使用once.Do
func GetConfig() *config {
	if cfg == nil {
		mux.Lock()
		defer mux.Unlock()
		if cfg == nil {
			vp := viper.New()
			vp.SetConfigName("mysql")
			vp.SetConfigType("yaml")
			vp.AddConfigPath("../../data/config")

			fmt.Println("Reading config file...")
			if err := vp.ReadInConfig(); err != nil {
				fmt.Println("Error reading config file:", err)
				return nil
			} else {
				// 解析配置文件
				cfg = &config{
					Password:      vp.GetString("lottery.pass"),
					ServerAddress: vp.GetString("lottery.host"),
				}
			}
		}
	}
	return cfg
}
