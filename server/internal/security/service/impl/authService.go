package impl

import (
	"SupCaller/internal/security/dto"
)

// AuthService 登录服务实现
type AuthService struct {
	// 可以在这里添加依赖，比如数据库连接等
}

// NewAuthService 创建登录服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login 实现 AuthServiceInterface 接口
func (s *AuthService) Login(loginForm *dto.LoginForm) (string, error) {
	return "12312", nil
}
