package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// 配置映射成Config对象

// 从Toml格式的配置文件加载配置
func LoadConfigFromToml(filePath string) error {
	config = NewDefaultConfig()

	// 读取toml格式配置
	_, err := toml.DecodeFile(filePath, config)
	if err != nil {
		return fmt.Errorf("load config from file error, path:%s, %s", filePath, err)
	}
	return nil
}

// 从环境变量加载配置
func LoadConfigFromEnv(filePath string) error {
	// 首先加载 .env 文件
	err := godotenv.Load(filePath)
	if err != nil {
		return fmt.Errorf("Error loading .env file: %v", err)
	}

	config = NewDefaultConfig()
	err = env.Parse(config)
	if err != nil {
		return err
	}
	return nil
}

// 加载全局实例
func loadGloabal() (err error) {
	// 加载db的全局实例
	db, err = config.MySQL.getDBConn()
	if err != nil {
		return
	}

	return
}
