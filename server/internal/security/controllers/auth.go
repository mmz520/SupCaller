package controllers

import (
	"SupCaller/common/response"
	"SupCaller/common/router"
	"SupCaller/internal/security/dto"
	"SupCaller/internal/security/service"
	"net/http"

	"SupCaller/internal/security/service/impl"

	"github.com/gin-gonic/gin"
)

// AuthService 登录服务实例
var authService service.AuthServiceInterface = impl.NewAuthService()

func init() {
	routerRegister := router.NewRouterRegister()
	routerRegister.Register("/auth", func(g *gin.RouterGroup) {
		g.POST("/login", LoginHandler)
		g.POST("/register", RegisterHandler)
	})
}

// LoginHandler 登录处理函数
func LoginHandler(c *gin.Context) {
	var loginForm dto.LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(500, err.Error()))
		return
	}
	data, err := authService.Login(&loginForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(data))
}

// RegisterHandler 注册处理函数
func RegisterHandler(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(nil))
}

// OAuth2AuthorizeHandler OAuth2授权处理函数
func OAuth2AuthorizeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(nil))
}

// OAuth2TokenHandler OAuth2令牌处理函数
func OAuth2TokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(nil))
}
