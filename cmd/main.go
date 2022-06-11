package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"gariwallet/api/proto/health/healthpb"
	"gariwallet/api/proto/wallet/walletpb"
	"gariwallet/api/proto/transaction/transactionpb"
	"gariwallet/internal/application/usecase"
	"gariwallet/internal/infrastructure"
	"gariwallet/internal/externalapp"
	"gariwallet/internal/presentation/middleware"
	"gariwallet/internal/presentation/server"
	"gariwallet/pkg/config"
	"gariwallet/pkg/database"
	"gariwallet/pkg/myjwt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type ErrorBody struct {
	Msg  string `json:"msg"`
	Code int32  `json:"code"`
}

func init() {
	log.SetPrefix("Wallet Server: ")
	if err := config.Load(); err != nil {
		panic(err)
	}
}

func main() {
	go func() {
		if err := RunGrpc(); err != nil {
			log.Fatalf("%s\n", err)
		}
	}()
	go func() {
		if err := RunGateway(); err != nil {
			log.Fatalf("%s\n", err)
		}
	}()
	select {}
}

func RunGrpc() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient := database.ConnectMongoDB(ctx)
	defer database.DisconnectDB(ctx, mongoClient)

	duration := time.Duration(config.Global.TokenDuration) * time.Minute
	jwtManager := myjwt.NewJwtManager(config.Global.JWTSignature, duration, duration+3600)
	authInterceptor := middleware.NewAuthInterceptor(jwtManager)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.AuthUnary()))

	wr := infrastructure.NewWalletRepository(mongoClient)
	bcs := externalapp.NewBlockchainServer()

	wa := usecase.NewWalletApp(wr)
	ta := usecase.NewTransactionApp(wr, bcs)
	walletServer := server.NewWalletManagementServer(wa)
	transactionServer := server.NewTransactionManagementServer(ta)
	healthServer := server.NewHealthServer()
	walletpb.RegisterWalletManagementServer(grpcServer, walletServer)
	transactionpb.RegisterTransactionManagementServer(grpcServer, transactionServer)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	reflection.Register(grpcServer)
	grpcAddress := fmt.Sprintf(":%s", config.Global.GrpcPort)
	grpcListener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("can't run grpc server")
	}

	return grpcServer.Serve(grpcListener)
}

func RunGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	runtime.HTTPError = CustomHTTPError
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := fmt.Sprintf("localhost:%s", config.Global.GrpcPort)
	err := walletpb.RegisterWalletManagementHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}
	err = transactionpb.RegisterTransactionManagementHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}
	err = healthpb.RegisterHealthHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(fmt.Sprintf(":%s", config.Global.RestPort), mux)
}

func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	body := &ErrorBody{
		Msg:  s.Message(),
		Code: int32(s.Code()),
	}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	// convert grpc code to http code
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}
