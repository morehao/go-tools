package jwtAuth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	type CustomData struct {
		Role string `json:"role"`
	}
	signKey := "secret"
	// uuid := uuid.NewString()
	uuid := "123456"
	now := time.Now()
	expiresAt := time.Now().Add(24 * time.Hour)
	issuedAt := time.Now()

	claims := NewClaims(
		WithCustomData(CustomData{Role: "admin"}),
		WithIssuer("example.com"),
		WithSubject("user123"),
		WithAudience("audience1", "audience2"),
		WithNotBefore(now),
		WithExpiresAt(expiresAt),
		WithIssuedAt(issuedAt),
		WithID(uuid),
	)
	token, err := CreateToken(signKey, claims)
	assert.Nil(t, err)
	t.Log(token)
}

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsInN1YiI6InVzZXIxMjMiLCJhdWQiOlsiYXVkaWVuY2UxIiwiYXVkaWVuY2UyIl0sImV4cCI6MTcyMzgwOTY4NCwibmJmIjoxNzIzNzIzMjg0LCJpYXQiOjE3MjM3MjMyODQsImp0aSI6IjEyMzQ1NiIsImN1c3RvbURhdGEiOnsicm9sZSI6ImFkbWluIn19.9a3KdeiA3Z9fK1pi2NrE-1nM3BVC4DdBY57GfGaCuts"
	signKey := "secret"
	type CustomerData struct {
		CompanyId uint64 `json:"companyId"`
		Role      string `json:"role"`
	}
	type CustomerClaims struct {
		CustomerData CustomerData `json:"customData"`
		jwt.RegisteredClaims
	}
	var claims CustomerClaims
	err := ParseToken(signKey, token, &claims)
	assert.Nil(t, err)
	t.Log(gutils.ToJsonString(claims))
	t.Log(claims.CustomerData.Role)
}

func TestRenewToken(t *testing.T) {
	signKey := "secret"
	newExpirationTime := 2 * time.Hour
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsInN1YiI6InVzZXIxMjMiLCJhdWQiOlsiYXVkaWVuY2UxIiwiYXVkaWVuY2UyIl0sImV4cCI6MTcyMzgwOTY4NCwibmJmIjoxNzIzNzIzMjg0LCJpYXQiOjE3MjM3MjMyODQsImp0aSI6IjEyMzQ1NiIsImN1c3RvbURhdGEiOnsicm9sZSI6ImFkbWluIn19.9a3KdeiA3Z9fK1pi2NrE-1nM3BVC4DdBY57GfGaCuts"
	newToken, err := RenewToken(signKey, token, newExpirationTime)
	assert.Nil(t, err)
	t.Log(newToken)
}
