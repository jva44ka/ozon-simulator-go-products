//go:generate swag init -g cmd/server/main.go --dir ./internal,./cmd
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jva44ka/ozon-simulator-go-products/internal/app/"
	"github.com/jva44ka/ozon-simulator-go-products/internal/middleware"
	desc "github.com/jva44ka/ozon-simulator-go-products/internal/pkg/api/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = ":50051"
	httpPort = ":8081"
)

func main() {
	ctx := context.Background()

	// запускаем gRPC сервер
	lis, _ := net.Listen("tcp", grpcPort)
	grpcServer := grpc.NewServer()
	pb.RegisterProductsServiceServer(grpcServer, server)

	go grpcServer.Serve(lis)

	// подключаемся к нему же как клиент
	conn, _ := grpc.Dial(httpPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// HTTP-прокси
	gwmux := runtime.NewServeMux()
	pb.RegisterProductsServiceHandler(ctx, gwmux, conn)

	http.ListenAndServe(":8080", gwmux)
}
