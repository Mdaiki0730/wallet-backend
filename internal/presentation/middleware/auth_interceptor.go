package middleware

import (
	"context"
	"errors"
	"log"

	"gariwallet/pkg/md"
	"gariwallet/pkg/myjwt"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
)

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	JwtManager *myjwt.JwtManager
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(jwtManager *myjwt.JwtManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

// handle authorization token
func (interceptor *AuthInterceptor) AuthUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		accessToken := md.GetTokenFromContext(ctx)

		parsedToken, err := interceptor.JwtManager.VerifyAuthorization(accessToken)
		if err != nil {
			log.Printf("token is unautherized\n")
			return nil, err
		}

		ctx, err = AddUserInfoToContext(ctx, parsedToken)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func AddUserInfoToContext(ctx context.Context, parsedToken *jwt.Token) (context.Context, error) {
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return ctx, errors.New("claims not found")
	}

	idpId, ok := claims["sub"].(string)
	if !ok {
		return ctx, errors.New("can't find user id")
	}
	ctx = context.WithValue(ctx, "idp_id", idpId)
	return ctx, nil
}
