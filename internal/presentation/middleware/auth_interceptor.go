package middleware

import (
	"context"
	"log"

	"gariwallet/pkg/md"
	"gariwallet/pkg/myjwt"

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

		_, err := interceptor.JwtManager.VerifyAuthorization(accessToken)
		if err != nil {
			log.Printf("token is unautherized\n")
			return nil, err
		}

		return handler(ctx, req)
	}
}
