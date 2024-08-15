package jwtAuth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
	"time"
)

type TokenOption struct {
}

func CreateToken(signKey string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signKey))
}

func ParseToken(signKey, token string, dest any) error {

	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Pointer {
		return errors.New("dest must be a pointer to a struct")
	}

	if destType.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	claims, ok := dest.(jwt.Claims)
	if !ok {
		return errors.New("dest does not implement jwt.Claims interface")
	}

	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	}
	tokenInst, err := jwt.ParseWithClaims(token, claims, keyFunc)

	if err != nil {
		return err
	}

	if !tokenInst.Valid {
		return errors.New("invalid token")
	}

	return nil
}

// RenewToken 续期 JWT
func RenewToken(signKey, oldToken string, newExpirationTime time.Duration) (string, error) {
	// 解析并验证旧的 token
	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	}
	token, err := jwt.ParseWithClaims(oldToken, &CustomClaims{}, keyFunc)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// 检查 token 是否有效
	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	// 获取旧的 claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", fmt.Errorf("cannot get claims from token")
	}

	// 更新过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(newExpirationTime))

	// 创建新的 token
	newTokenInst := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newTokenInst.SignedString([]byte(signKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign new token: %w", err)
	}

	return newTokenString, nil
}
