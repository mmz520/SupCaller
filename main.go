package main

import (
	"SupCaller/common/config"
	"SupCaller/common/database"
	"SupCaller/internal/security/controllers"
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
	{
		// 认证相关路由组
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.LoginHandler)
			auth.POST("/register", controllers.RegisterHandler)
			auth.POST("/oauth2/authorize", controllers.OAuth2AuthorizeHandler)
			auth.POST("/oauth2/token", controllers.OAuth2TokenHandler)
		}
	}
}

func listUsers(c *gin.Context) {
	// 实际实现用户列表逻辑
	c.JSON(200, gin.H{"message": "List users endpoint"})
}

func createUser(c *gin.Context) {
	// 实际实现创建用户逻辑
	c.JSON(200, gin.H{"message": "Create user endpoint"})
}

func getUser(c *gin.Context) {
	// 实际实现获取用户逻辑
	c.JSON(200, gin.H{"message": "Get user endpoint"})
}

func updateUser(c *gin.Context) {
	// 实际实现更新用户逻辑
	c.JSON(200, gin.H{"message": "Update user endpoint"})
}

func deleteUser(c *gin.Context) {
	// 实际实现删除用户逻辑
	c.JSON(200, gin.H{"message": "Delete user endpoint"})
}

func listRoles(c *gin.Context) {
	// 实际实现角色列表逻辑
	c.JSON(200, gin.H{"message": "List roles endpoint"})
}

func createRole(c *gin.Context) {
	// 实际实现创建角色逻辑
	c.JSON(200, gin.H{"message": "Create role endpoint"})
}

func getRole(c *gin.Context) {
	// 实际实现获取角色逻辑
	c.JSON(200, gin.H{"message": "Get role endpoint"})
}

func updateRole(c *gin.Context) {
	// 实际实现更新角色逻辑
	c.JSON(200, gin.H{"message": "Update role endpoint"})
}

func deleteRole(c *gin.Context) {
	// 实际实现删除角色逻辑
	c.JSON(200, gin.H{"message": "Delete role endpoint"})
}

func listPermissions(c *gin.Context) {
	// 实际实现权限列表逻辑
	c.JSON(200, gin.H{"message": "List permissions endpoint"})
}

func createPermission(c *gin.Context) {
	// 实际实现创建权限逻辑
	c.JSON(200, gin.H{"message": "Create permission endpoint"})
}

func getPermission(c *gin.Context) {
	// 实际实现获取权限逻辑
	c.JSON(200, gin.H{"message": "Get permission endpoint"})
}

func updatePermission(c *gin.Context) {
	// 实际实现更新权限逻辑
	c.JSON(200, gin.H{"message": "Update permission endpoint"})
}

func deletePermission(c *gin.Context) {
	// 实际实现删除权限逻辑
	c.JSON(200, gin.H{"message": "Delete permission endpoint"})
}
