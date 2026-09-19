package database

import (
	"fmt"
	"wk-go/internal/config"
	"wk-go/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	cfg := config.Global.Database
	return InitWithDSN(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset,
	))
}

// InitWithDSN 用指定 DSN 初始化数据库连接（安装引导用）
func InitWithDSN(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return err
	}
	DB = db

	// 自动建表/补字段
	db.AutoMigrate(
		&model.User{},
		&model.HuoYuan{},
		&model.Order{},
		&model.Class{},
		&model.FenLei{},
	)

	return nil
}
