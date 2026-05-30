package rbac

import (
	"SupCaller/common/config"
	"SupCaller/internal/security/models"
	"strings"

	"gorm.io/gorm"
)

// RBAC 权限控制结构体
type RBAC struct {
	db *gorm.DB
}

// New 创建RBAC实例
func New(db *gorm.DB) *RBAC {
	return &RBAC{
		db: db,
	}
}

// GetUserRoles 获取用户角色
func (r *RBAC) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Table("roles r").
		Joins("left join user_roles ur on r.id = ur.role_id").
		Where("ur.user_id = ? and r.status = 1", userID).
		Find(&roles).Error
	return roles, err
}

// GetRolePermissions 获取角色权限
func (r *RBAC) GetRoleAuth(roleID uint) ([]models.Auth, error) {
	var auth []models.Auth
	err := r.db.Table("auth p").
		Joins("left join role_auth rp on p.id = rp.auth_id").
		Where("rp.role_id = ? and p.status = 1", roleID).
		Find(&auth).Error
	return auth, err
}

// GetUserPermissions 获取用户所有权限
func (r *RBAC) GetUserAuth(userID uint) ([]models.Auth, error) {
	var auth []models.Auth
	err := r.db.Raw(`
		SELECT DISTINCT p.* 
		FROM auth p
		LEFT JOIN role_auth rp ON p.id = rp.auth_id
		LEFT JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND p.status = 1
	`, userID).Scan(&auth).Error
	return auth, err
}

// CheckAuth 检查用户是否具有指定权限
func (r *RBAC) CheckAuth(userID uint, authCode string) (bool, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) 
		FROM auth p
		LEFT JOIN role_auth rp ON p.id = rp.auth_id
		LEFT JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND p.code = ? AND p.status = 1
	`, userID, authCode).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckInterfaceAuth 检查用户是否具有接口权限
func (r *RBAC) CheckInterfaceAuth(userID uint, interfacePath string) (bool, error) {
	// 检查是否是忽略的接口
	for _, ignorePath := range config.Config.Ignore {
		if matchPath(interfacePath, ignorePath) {
			return true, nil
		}
	}
	// 将接口路径转换为权限代码格式
	authCode := strings.ReplaceAll(interfacePath, "/", ":")
	authCode = strings.TrimPrefix(authCode, ":")
	return r.CheckAuth(userID, authCode)
}

// matchPath 检查接口路径是否匹配忽略模式
// 支持的模式:
// - /path/path2 - 精确匹配
// - /path/** - 匹配 /path/ 下的所有路径（包括 /path/ 本身）
// - /path/**/path2 - 匹配 /path/ 和 /path2 之间有任意层级的路径
func matchPath(path, pattern string) bool {
	// 精确匹配
	if path == pattern {
		return true
	}
	
	// 如果模式不包含 **，则不匹配
	if !strings.Contains(pattern, "**") {
		return false
	}
	
	// 分割路径和模式
	pathParts := strings.Split(path, "/")
	patternParts := strings.Split(pattern, "/")
	
	return matchParts(pathParts, patternParts)
}

// matchParts 递归匹配路径片段
func matchParts(pathParts, patternParts []string) bool {
	for len(patternParts) > 0 {
		part := patternParts[0]
		
		if part == "**" {
			// 移除 **
			patternParts = patternParts[1:]
			
			// 如果 ** 是最后一部分，匹配成功
			if len(patternParts) == 0 {
				return true
			}
			
			// ** 可以匹配零个或多个路径段，尝试所有可能性
			for i := 0; i <= len(pathParts); i++ {
				if matchParts(pathParts[i:], patternParts) {
					return true
				}
			}
			return false
		}
		
		// 如果路径已用完但模式还有剩余，不匹配
		if len(pathParts) == 0 {
			return false
		}
		
		// 当前段不匹配
		if pathParts[0] != part {
			return false
		}
		
		// 继续匹配下一段
		pathParts = pathParts[1:]
		patternParts = patternParts[1:]
	}
	
	// 模式匹配完，检查路径是否也用完了
	return len(pathParts) == 0
}
