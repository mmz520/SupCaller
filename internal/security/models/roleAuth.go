package models

type RoleAuth struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	RoleID uint `gorm:"index" json:"role_id"`
	AuthID uint `gorm:"index" json:"auth_id"`
}
