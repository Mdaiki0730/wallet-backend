package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"gariwallet/api/proto/wallet/walletpb"
	"gariwallet/internal/application/usecase"
	"gariwallet/internal/infrastructure"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient := database.ConnectMongoDB(ctx)
	defer database.DisconnectDB(ctx, mongoClient)

	duration := time.Duration(config.Global.TokenDuration) * time.Minute
	jwtManager := myjwt.NewJwtManager(config.Global.JWTSignature, duration, duration+3600)
	authInterceptor := middleware.NewAuthInterceptor(jwtManager)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.AuthUnary()))

	wr := infrastructure.NewWalletRepository(mongoClient)
	wa := usecase.NewWalletApp(wr)
	walletServer := server.NewWalletManagementServer(wa)
	walletpb.RegisterWalletManagementServer(grpcServer, walletServer)

	reflection.Register(grpcServer)
	grpcAddress := fmt.Sprintf(":%s", config.Global.GrpcPort)
	grpcListener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("can't run grpc server")
	}

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
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
