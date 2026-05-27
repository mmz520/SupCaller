package models

import "time"

type Auth struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100" json:"name"`
	Code        string    `gorm:"size:100;uniqueIndex" json:"code"` // 权限代码，如：user:create, user:delete
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
