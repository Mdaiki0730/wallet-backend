package md

import (
	"context"
	"strings"

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

	var jwtString string
	if strings.Contains(authorization[0], "Bearer ") {
		jwtStrings := strings.Split(authorization[0], "Bearer ")
		if len(jwtStrings) > 1 {
			jwtString = jwtStrings[1]
		}
	} else {
		jwtString = authorization[0]
	}

	return jwtString
}
