package server

import (
	"context"

	"gariwallet/api/proto/wallet/walletpb"
	"gariwallet/internal/application/usecase"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

type walletManagementServer struct {
	application usecase.WalletAppIF
}

func NewWalletManagementServer(app usecase.WalletAppIF) walletpb.WalletManagementServer {
	return &walletManagementServer{app}
}

func (wm *walletManagementServer) Create(ctx context.Context, req *walletpb.CreateWalletRequest) (*walletpb.WalletBaseResponse, error) {
	result, err := wm.application.Create(ctx)
	if err != nil {
		return nil, err
	}

	return &walletpb.WalletBaseResponse{BlockchainAddress: &result.BlockchainAddress}, nil
}
