package main

import (
	"SupCaller/common/config"
	"SupCaller/common/database"
	"SupCaller/common/router"
	"SupCaller/internal/security/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	database.InitDB()

	// 创建Gin引擎
	r := gin.Default()

	// 初始化认证中间件
	authMiddleware := middleware.NewAuthMiddleware(
		database.DB,
		middleware.WithIgnore(config.Config.Ignore),
	)

	// 使用认证中间件
	r.Use(authMiddleware.Handle())

	// 注册路由
	setupRoutes(r)

	// 启动服务
	addr := fmt.Sprintf(":%d", config.Config.App.Port)
	log.Printf("Server starting on port %d", config.Config.App.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "SupCaller service is running",
		})
	})
	// API路由组
	api := r.Group("/api/v1")
	router.AutoRegister(api)
}
