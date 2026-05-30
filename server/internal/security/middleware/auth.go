package middleware

import (
	"SupCaller/common/response"
	"SupCaller/common/utils"
	"SupCaller/internal/security/rbac"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware 认证中间件结构体
type AuthMiddleware struct {
	db     *gorm.DB
	rbac   *rbac.RBAC
	ignore []string
}

// AuthOption 认证中间件选项
type AuthOption func(*AuthMiddleware)

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(db *gorm.DB, opts ...AuthOption) *AuthMiddleware {
	am := &AuthMiddleware{
		db:   db,
		rbac: rbac.New(db),
		ignore: []string{
			"/auth/login",
			"/auth/register",
			"/health",
			"/public/**",
			"/assets/**",
		},
	}

	for _, opt := range opts {
		opt(am)
	}

	return am
}

// WithIgnore 设置忽略路径
func WithIgnore(paths []string) AuthOption {
	return func(am *AuthMiddleware) {
		am.ignore = paths
	}
}

// isIgnore 检查路径是否在忽略列表中
func (am *AuthMiddleware) isIgnore(path string) bool {
	for _, pattern := range am.ignore {
		// 检查精确匹配
		if path == pattern {
			return true
		}

		// 检查通配符匹配 /**
		if strings.HasSuffix(pattern, "/**") {
			prefix := pattern[:len(pattern)-3] // 去掉 /**
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}

	return false
}

// Handle 认证中间件处理函数
func (am *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在忽略列表中
		if am.isIgnore(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error(401, "缺少访问令牌"))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, response.Error(401, "令牌格式错误"))
			c.Abort()
			return
		}

		// 提取token
		tokenString := authHeader[7:] // 去掉 "Bearer " 前缀

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error(401, err.Error()))
			c.Abort()
			return
		}
		// 将用户信息存储到上下文中
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)

		// 检查接口权限
		hasPermission, err := am.rbac.CheckInterfaceAuth(claims.UserID, c.Request.URL.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, "权限检查失败"))
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, response.Error(403, "无访问权限"))
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}
