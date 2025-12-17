package secure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var (
	jwtKey []byte
	expiresTime time.Duration
	refreshTime time.Duration
)

func InitJWT() error {
	key := viper.GetString("jwt.key")
	if key == "" {
		return errors.New("jwtKey是空的")
	}

	jwtKey = []byte(key)
	expiresTime = viper.GetDuration("jwt.expires_time")
	refreshTime = viper.GetDuration("jwt.refresh_time")
	return nil
}

func GetExpiresTime() time.Duration {
	return expiresTime
}

func GetRefreshtime() time.Duration {
	return refreshTime
}

type IDClaims struct {
	ID uint64 `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(id uint64, t time.Duration, token_type string) (string, error) {
	if string(jwtKey) == "" {
		return "", errors.New("jwtKey在生成Token时发现是空的")
	}

	claims := IDClaims{
		ID: id,
		Type: token_type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并生成字符串
	return token.SignedString(jwtKey)
}

func NewToken(id uint64) (string, error) {
	return GenerateToken(id, expiresTime, "access")
}

func NewRefreshToken(id uint64) (string, error) {
	return GenerateToken(id, refreshTime, "refresh")
}
