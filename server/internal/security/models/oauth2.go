package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"gorm.io/gorm"
)

// OAuth2Client OAuth2客户端模型
type OAuth2Client struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ClientID     string    `gorm:"uniqueIndex;size:100" json:"client_id"`
	ClientSecret string    `gorm:"size:255" json:"client_secret"`
	Name         string    `gorm:"size:100" json:"name"`
	RedirectURL  string    `gorm:"size:255" json:"redirect_url"`
	Status       int       `gorm:"default:1" json:"status"` // 1:正常 2:禁用
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OAuth2AuthorizationCode OAuth2授权码模型
type OAuth2AuthorizationCode struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"index;size:255" json:"code"`
	ClientID    string    `gorm:"size:100" json:"client_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	RedirectURL string    `gorm:"size:255" json:"redirect_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// OAuth2AccessToken OAuth2访问令牌模型
type OAuth2AccessToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"index;size:500" json:"token"`
	ClientID  string    `gorm:"size:100" json:"client_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

// CreateClient 创建OAuth2客户端
func CreateClient(db *gorm.DB, name, redirectURL string) (*OAuth2Client, error) {
	clientID, err := GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	clientSecret, err := GenerateRandomString(48)
	if err != nil {
		return nil, err
	}

	client := &OAuth2Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Name:         name,
		RedirectURL:  redirectURL,
		Status:       1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(client).Error; err != nil {
		return nil, err
	}

	return client, nil
}

// GenerateAuthorizationCode 生成授权码
func GenerateAuthorizationCode(db *gorm.DB, clientID string, userID uint, redirectURL string) (*OAuth2AuthorizationCode, error) {
	code, err := GenerateRandomString(64)
	if err != nil {
		return nil, err
	}

	authCode := &OAuth2AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      userID,
		RedirectURL: redirectURL,
		ExpiresAt:   time.Now().Add(5 * time.Minute), // 5分钟过期
		CreatedAt:   time.Now(),
	}

	if err := db.Create(authCode).Error; err != nil {
		return nil, err
	}

	return authCode, nil
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(db *gorm.DB, clientID string, userID uint) (*OAuth2AccessToken, error) {
	token, err := GenerateRandomString(128)
	if err != nil {
		return nil, err
	}

	accessToken := &OAuth2AccessToken{
		Token:     token,
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(2 * time.Hour), // 2小时过期
		CreatedAt: time.Now(),
	}

	if err := db.Create(accessToken).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}
