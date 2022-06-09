package md

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetTokenFromContext(ctx context.Context) (token string) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return token
	}

	authorization := md["authorization"]
	if len(authorization) == 0 {
		return token
	}

	return authorization[0]
}
