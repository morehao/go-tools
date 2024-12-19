package jwtAuth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 包含注册的 claims 和自定义的 claims
type Claims struct {
	jwt.RegisteredClaims
	CustomData any `json:"customData,omitempty"` // 自定义的结构
}

// ClaimsOption 定义用于配置 jwt.RegisteredClaims 和自定义 Claims 的函数类型
type ClaimsOption func(*Claims)

// WithIssuer 配置 Issuer 声明
func WithIssuer(issuer string) ClaimsOption {
	return func(c *Claims) {
		c.Issuer = issuer
	}
}

// WithSubject 配置 Subject 声明
func WithSubject(subject string) ClaimsOption {
	return func(c *Claims) {
		c.Subject = subject
	}
}

// WithAudience 配置 Audience 声明
func WithAudience(audience ...string) ClaimsOption {
	return func(c *Claims) {
		c.Audience = audience
	}
}

// WithExpiresAt 配置 ExpiresAt 声明
func WithExpiresAt(expiresAt time.Time) ClaimsOption {
	return func(c *Claims) {
		c.ExpiresAt = jwt.NewNumericDate(expiresAt)
	}
}

// WithNotBefore 配置 NotBefore 声明
func WithNotBefore(notBefore time.Time) ClaimsOption {
	return func(c *Claims) {
		c.NotBefore = jwt.NewNumericDate(notBefore)
	}
}

// WithIssuedAt 配置 IssuedAt 声明
func WithIssuedAt(issuedAt time.Time) ClaimsOption {
	return func(c *Claims) {
		c.IssuedAt = jwt.NewNumericDate(issuedAt)
	}
}

// WithID 配置 ID 声明
func WithID(id string) ClaimsOption {
	return func(c *Claims) {
		c.ID = id
	}
}

// WithCustomData 配置自定义 Claims
func WithCustomData(custom interface{}) ClaimsOption {
	return func(c *Claims) {
		c.CustomData = custom
	}
}

// NewClaims 创建并配置 Claims 实例
func NewClaims(opts ...ClaimsOption) *Claims {
	claims := &Claims{}
	for _, opt := range opts {
		opt(claims)
	}
	return claims
}
