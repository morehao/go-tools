package jwtAuth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	type CustomerData struct {
		CompanyId uint64
	}
	signKey := "secret"
	// uuid := uuid.NewString()
	uuid := "123456"
	claims := &Claims{
		CustomData: CustomerData{CompanyId: 1},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "test",
			Subject:   "test",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid,
		},
	}
	token, err := CreateToken(signKey, claims)
	assert.Nil(t, err)
	t.Log(token)
}

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0Iiwic3ViIjoidGVzdCIsImV4cCI6MTcyMzgwOTA2OCwiaWF0IjoxNzIzNzIyNjY4LCJqdGkiOiIxMjM0NTYiLCJjdXN0b21EYXRhIjp7IkNvbXBhbnlJZCI6MX19.isZMExv6HbQYmQuYMKZ1sgVcCmLzBFswXbMJKY1ibP8"
	signKey := "secret"
	type CustomerData struct {
		CompanyId uint64
	}
	type CustomerClaims struct {
		CustomerData CustomerData `json:"customData"`
		jwt.RegisteredClaims
	}
	var claims CustomerClaims
	err := ParseToken(signKey, token, &claims)
	assert.Nil(t, err)
	t.Log(gutils.ToJsonString(claims))
	t.Log(claims.CustomerData.CompanyId)
}

func TestRenewToken(t *testing.T) {
	// 自定义 claims 结构体
	type MyCustomClaims struct {
		Role string `json:"role"`
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	claims := NewClaims(
		MyCustomClaims{Role: "admin"},
		WithIssuer("example.com"),
		WithSubject("user123"),
		WithAudience("audience1", "audience2"),
		WithNotBefore(now),
		WithExpiresAt(expiresAt),
		WithIssuedAt(now),
		WithID("unique-id-12345"),
	)
	signKey := "secret"
	token, err := CreateToken(signKey, claims)
	assert.Nil(t, err)
	newExpirationTime := 2 * time.Hour
	newToken, err := RenewToken(signKey, token, newExpirationTime)
	assert.Nil(t, err)
	t.Log(newToken)
}
