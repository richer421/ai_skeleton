package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load 加载配置
func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")
	viper.AddConfigPath("..")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}
