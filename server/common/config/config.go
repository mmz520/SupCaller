package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 全局配置变量
var Config *Configuration

// Configuration 配置结构体
type Configuration struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	OAuth2   OAuth2Config   `mapstructure:"oauth2"`
	Ignore   []string       `mapstructure:"ignore"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Env       string `mapstructure:"env"`
	Port      int    `mapstructure:"port"`
	SecretKey string `mapstructure:"secret_key"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	DBName    string `mapstructure:"dbname"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parseTime"`
	Loc       string `mapstructure:"loc"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// OAuth2Config OAuth2配置
type OAuth2Config struct {
	AccessTokenExpire  string `mapstructure:"access_token_expire"`
	RefreshTokenExpire string `mapstructure:"refresh_token_expire"`
	AuthCodeExpire     string `mapstructure:"authorization_code_expire"`
}

// LoadConfig 加载配置文件
func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	// 设置默认值
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.env", "development")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 解析配置到结构体
	Config = &Configuration{}
	if err := viper.Unmarshal(Config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
