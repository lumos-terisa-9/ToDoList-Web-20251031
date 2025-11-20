package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	DebugMode  bool   `mapstructure:"DEBUG"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

func LoadConfig() *Config {
	// 设置viper配置
	viper.SetConfigName(".env") // 配置文件名称
	viper.SetConfigType("env")  // 配置文件类型
	viper.AddConfigPath(".")    //找到env文件所在目录

	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("警告：无法读取配置文件: %v", err)
		log.Println("将使用环境变量和默认值")
	}

	// 自动绑定环境变量
	viper.AutomaticEnv()

	// 将配置解析到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	log.Printf("配置加载成功 - 服务器端口: %s", config.ServerPort)
	log.Printf("Redis配置: %s, DB: %d", config.RedisAddr, config.RedisDB)

	return &config
}

func setDefaults() {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DEBUG", true)

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "team_task_hub")

	viper.SetDefault("JWT_SECRET", "default-jwt-secret-change-in-production")

	viper.SetDefault("SMTP_HOST", "smtp.qq.com")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USER", "")
	viper.SetDefault("SMTP_PASSWORD", "")
	viper.SetDefault("SMTP_FROM", "")

	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
}
