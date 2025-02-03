package main

import (
	"fmt"
	"gym/internal/handler"
	"gym/internal/middleware"
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

	// 创建处理器
	loginHandler := handler.NewLoginHandler()
	userHandler := handler.NewUserHandler()

	// 注册路由
	api := r.Group("/api/v1")
	{
		// 公开路由 - 登录相关
		login := api.Group("/login")
		{
			login.POST("/wx-mini", loginHandler.WxMiniLogin) // 微信小程序登录
			login.POST("/phone", loginHandler.PhoneLogin)    // 手机号密码登录
			login.POST("/sms", loginHandler.SmsLogin)        // 短信验证码登录
		}

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWT())
		{
			// 用户相关
			auth.GET("/user/profile", userHandler.GetProfile)       // 获取用户信息
			auth.PUT("/user/profile", userHandler.UpdateProfile)    // 更新用户信息
			auth.POST("/user/password", userHandler.ChangePassword) // 修改密码
		}
	}

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
