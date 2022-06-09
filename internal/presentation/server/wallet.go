package server

import (
	"context"

	"gariwallet/api/proto/wallet/walletpb"
	"gariwallet/internal/application/command"
	"gariwallet/internal/application/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type walletManagementServer struct {
	application usecase.WalletAppIF
}

func NewWalletManagementServer(app usecase.WalletAppIF) walletpb.WalletManagementServer {
	return &walletManagementServer{app}
}

func (wm *walletManagementServer) Create(ctx context.Context, req *walletpb.CreateWalletRequest) (*walletpb.WalletBaseResponse, error) {
	cmd := command.WalletCreate{ctx.Value("idp_id").(string)}
	result, err := wm.application.Create(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &walletpb.WalletBaseResponse{BlockchainAddress: &result.BlockchainAddress}, nil
}

func (wm *walletManagementServer) Delete(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	cmd := command.WalletDelete{ctx.Value("idp_id").(string)}
	err := wm.application.Delete(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (wm *walletManagementServer) Get(ctx context.Context, req *emptypb.Empty) (*walletpb.WalletBaseResponse, error) {
	cmd := command.WalletGet{ctx.Value("idp_id").(string)}
	result, err := wm.application.Get(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &walletpb.WalletBaseResponse{BlockchainAddress: &result.BlockchainAddress}, nil
}
