package untils

import (
	"SupCaller/common/config"

	"SupCaller/internal/security/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成token的方法
func GenerateToken(user models.User) string {
	nowTime := time.Now()
	//config.Config.JWT.Expire是24h, 转换为time.Duration
	millisecond, err := time.ParseDuration(config.Config.JWT.Expire)
	if err != nil {
		panic(errors.New("config.jwt.expire time format is invalid"))
	}
	expireTime := nowTime.Add(millisecond * time.Second)
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "supcaller",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JWT.Secret))
	if err != nil {
		panic(errors.New("token sign error"))
	}
	return tokenString
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWT.Secret), nil
	})
	if err != nil {
		// 检查是否是因为token过期导致的错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, err
		}
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, err
}
