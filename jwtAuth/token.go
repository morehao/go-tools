package jwtAuth

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(signKey string, claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signKey))
}

func ParseToken(signKey, tokenStr string, dest any) error {
	// 检查 dest 是否为指向结构体的指针
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Pointer || destType.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	// 检查 dest 是否实现了 jwt.Claims 接口
	claims, ok := dest.(jwt.Claims)
	if !ok {
		return errors.New("dest does not implement jwt.Claims interface")
	}

	// 定义用于解析 JWT 的 keyFunc
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	}

	// 解析 JWT
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		return err
	}

	// 检查 token 是否有效
	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

// RenewToken 续期 JWT
func RenewToken(signKey, oldTokenStr string, newExpirationTime time.Duration) (string, error) {
	// 解析并验证旧的 token
	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	}
	token, err := jwt.ParseWithClaims(oldTokenStr, &Claims{}, keyFunc)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// 检查 token 是否有效
	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	// 获取旧的 claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("cannot get claims from token")
	}

	// 更新过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(newExpirationTime))

	// 创建新的 token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString([]byte(signKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign new token: %w", err)
	}

	return newTokenString, nil
}
