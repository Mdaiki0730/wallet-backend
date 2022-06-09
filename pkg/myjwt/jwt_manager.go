package myjwt

import (
	"errors"
	"strings"
	"time"

	"gariwallet/pkg/config"

	"github.com/golang-jwt/jwt"
)

// JwtManager is a JSON web token manager
type JwtManager struct {
	signature            string
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
}

// NewJwtManager returns a new JWT manager
func NewJwtManager(signature string, tokenDuration, refreshTokenDuration time.Duration) *JwtManager {
	return &JwtManager{signature, tokenDuration, refreshTokenDuration}
}

// VerifyAuthorization verifies the access token string
func (manager *JwtManager) VerifyAuthorization(authorization string) (*jwt.Token, error) {
	var jwtString string
	if strings.Contains(authorization, "Bearer ") {
		jwtStrings := strings.Split(authorization, "Bearer ")
		if len(jwtStrings) > 1 {
			jwtString = jwtStrings[1]
		}
	} else {
		jwtString = authorization
	}

	parsedToken, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("not match signing method")
		}
		cert := "-----BEGIN CERTIFICATE-----\n" + config.Global.JWTSignature + "\n-----END CERTIFICATE-----"
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, errors.New("generated invalid pem")
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New("token is invalid")
	}

	return parsedToken, nil
}
