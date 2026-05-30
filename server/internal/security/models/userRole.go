// UserRole 用户角色关联
package models

type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index" json:"user_id"`
	RoleID uint `gorm:"index" json:"role_id"`
}
