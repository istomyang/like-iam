package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	JwtExpireDevelopment = time.Hour * 24 * 7
	JwtExpireProduction  = time.Hour * 2
)

// Sign put secretId into header section and signs jwt-token with secretKey.
// issuer is signer name, audience is consumer of this token.
func Sign(secretId, secretKey, issuer, audience string, expire time.Duration) (string, error) {
	// https://datatracker.ietf.org/doc/html/rfc7519#section-4-1
	var token = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Issuer:    issuer,
		Audience:  []string{audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	token.Header["kid"] = secretId
	return token.SignedString(secretKey)
}
