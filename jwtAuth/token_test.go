package jwtAuth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	signKey := "secret"
	type CustomerClaims struct {
		CompanyId uint64
		jwt.RegisteredClaims
	}
	// uuid := uuid.NewString()
	uuid := "123456"
	claims := CustomerClaims{
		CompanyId: 1,
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
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDb21wYW55SWQiOjEsImlzcyI6InRlc3QiLCJzdWIiOiJ0ZXN0IiwiZXhwIjoxNzIzNzk1NDQ4LCJpYXQiOjE3MjM3MDkwNDh9.C_7qYBPr1HWeSRmQ3vPcbvZP3sNh1HRceFyeZ17vGkU"
	signKey := "secret"
	type CustomerClaims struct {
		CompanyId uint64
		jwt.RegisteredClaims
	}
	var claims CustomerClaims
	err := ParseToken(signKey, token, &claims)
	assert.Nil(t, err)
	t.Log(gutils.ToJsonString(claims))
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
