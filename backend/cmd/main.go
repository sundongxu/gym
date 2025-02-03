package main

import (
	"fmt"
	"gym/pkg/config"
	"gym/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", err))
	}

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		panic(fmt.Sprintf("初始化数据库连接失败: %v", err))
	}

	// 获取底层sqlDB连接，用于后续关闭
	sqlDB, err := database.DB.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库连接失败: %v", err))
	}
	defer sqlDB.Close()

	// 自动迁移数据库表结构
	if err := database.AutoMigrate(); err != nil {
		panic(fmt.Sprintf("数据库迁移失败: %v", err))
	}

	// 设置gin模式
	gin.SetMode(config.GetConfig().Server.Mode)

	// 创建gin引擎
	r := gin.Default()

	// 注册基本健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务正常运行",
		})
	})

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.GetConfig().Server.Port)
	if err := r.Run(addr); err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
}
