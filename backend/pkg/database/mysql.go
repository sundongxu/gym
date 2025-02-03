package database

import (
	"fmt"
	"gym/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	// ErrRecordNotFound record not found error
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func Init() error {
	var err error
	conf := config.GetConfig().MySQL

	// 先创建数据库
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Charset,
	)

	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to mysql failed: %v", err)
	}

	// 创建数据库
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", conf.DBName)
	err = rootDB.Exec(createDBSQL).Error
	if err != nil {
		return fmt.Errorf("create database failed: %v", err)
	}

	// 连接指定的数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.Charset,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to mysql failed: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %v", err)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
