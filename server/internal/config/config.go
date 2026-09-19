package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Admin    AdminConfig    `mapstructure:"admin"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

type AdminConfig struct {
	VerifyPassword string `mapstructure:"verify_password"`
}

var Global *Config

func Init(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	Global = &Config{}
	return viper.Unmarshal(Global)
}

// Save 把当前配置写回配置文件（安装引导用）
func Save() error {
	if Global == nil {
		return fmt.Errorf("配置未初始化")
	}
	viper.Set("server.port", Global.Server.Port)
	viper.Set("server.mode", Global.Server.Mode)
	viper.Set("database.host", Global.Database.Host)
	viper.Set("database.port", Global.Database.Port)
	viper.Set("database.user", Global.Database.User)
	viper.Set("database.password", Global.Database.Password)
	viper.Set("database.dbname", Global.Database.DBName)
	viper.Set("database.charset", Global.Database.Charset)
	viper.Set("jwt.secret", Global.JWT.Secret)
	viper.Set("jwt.expire", Global.JWT.Expire)
	viper.Set("admin.verify_password", Global.Admin.VerifyPassword)
	return viper.WriteConfig()
}
