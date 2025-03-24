package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
	AI struct {
		ApiKey string `mapstructure:"apiKey"`
		Url    string `mapstructure:"url"`
		Model  string `mapstructure:"model"`
		Prompt string `mapstructure:"prompt"`
	} `mapstructure:"ai"`
	Server struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		LogLevel string `mapstructure:"logLevel"`
	} `mapstructure:"server"`
}

var AppConfig Config

// LoadConfig 使用 Viper 加载配置
func LoadConfig() error {
	env := os.Getenv("BLOG_ENV")
	if env == "DEV" {
		viper.SetConfigName("config-dev") // 配置文件名称（不带扩展名）
	} else {
		viper.SetConfigName("config") // 配置文件名称（不带扩展名）

	}
	viper.SetConfigType("yaml")    // 配置文件类型
	viper.AddConfigPath("config/") // 配置文件路径

	// 设置环境变量前缀
	viper.SetEnvPrefix("BLOG")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 将配置映射到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	// 监听配置文件变化（可选）
	viper.WatchConfig()

	return nil
}

// GetDSN 获取数据库连接字符串
func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.Name)
}
