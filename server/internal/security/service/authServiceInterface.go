package service

import "SupCaller/internal/security/dto"

// AuthServiceInterface 登录服务接口
type AuthServiceInterface interface {
	Login(loginForm *dto.LoginForm) (string, error)
}
