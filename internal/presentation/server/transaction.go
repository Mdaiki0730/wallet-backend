package server

import (
	"context"
	"encoding/json"

	"gariwallet/api/proto/transaction/transactionpb"
	"gariwallet/internal/application/command"
	"gariwallet/internal/application/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type transactionManagementServer struct {
	application usecase.TransactionAppIF
}

func NewTransactionManagementServer(app usecase.TransactionAppIF) transactionpb.TransactionManagementServer {
	return &transactionManagementServer{app}
}

func (tm *transactionManagementServer) Create(ctx context.Context, req *transactionpb.CreateTransactionRequest) (*transactionpb.TransactionBaseResponse, error) {
	cmd := command.TransactionCreate{ctx.Value("idp_id").(string), req.GetRecipientBlockchainAddress(), req.GetValue()}
	result, err := tm.application.Create(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res := &transactionpb.TransactionBaseResponse{}
	b, _ := json.Marshal(result)
	json.Unmarshal(b, res)
	return res, nil
}
