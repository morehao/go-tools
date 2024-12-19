package jwtAuth

import (
	"testing"
	"time"

	"github.com/morehao/go-tools/gutils"
)

func TestNewClaims(t *testing.T) {
	// 自定义 claims 结构体
	type CustomData struct {
		Role string `json:"role"`
	}

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
		WithID("unique-id-12345"),
	)

	t.Log(gutils.ToJsonString(claims))
}
