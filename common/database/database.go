package database

import (
	"SupCaller/common/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.Config.Database.Username,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.DBName,
		config.Config.Database.Charset,
		config.Config.Database.ParseTime,
		config.Config.Database.Loc,
	)

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据表
	err = DB.AutoMigrate(
	// &model.User{},
	// &model.Role{},
	// &model.Permission{},
	// &model.UserRole{},
	// &model.RolePermission{},
	// &model.OAuth2Client{},
	// &model.OAuth2AuthorizationCode{},
	// &model.OAuth2AccessToken{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}
