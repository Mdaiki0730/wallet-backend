package server

import (
	"context"

	"gariwallet/api/proto/health/healthpb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type healthServer struct{}

func NewHealthServer() healthpb.HealthServer {
	return &healthServer{}
}

func (hs *healthServer) Check(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
