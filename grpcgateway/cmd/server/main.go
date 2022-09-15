package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	productRPC "grpcgateway/api/gRPC/product"
	"grpcgateway/internal/pb/product"
	"log"
	"net"
	"net/http"
	"sync"
)

const (
	grpcPort = ":8080"
	httpPort = ":8081"
)

func main() {
	grpcL, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serverGrpc := grpc.NewServer()
	product.RegisterProductInfoServer(serverGrpc, &productRPC.Server{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup, listener net.Listener) {
		defer wg.Done()
		if errGrpc := serverGrpc.Serve(listener); errGrpc != nil {
			log.Fatal("unable to start server", errGrpc)
		}
	}(&wg, grpcL)
	log.Printf("grpc server started")
	mux := runtime.NewServeMux()
	err = product.RegisterProductInfoHandlerFromEndpoint(context.Background(),
		mux,
		fmt.Sprintf("localhost%v", grpcPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatal(err)
	}
	serverHTTP := http.Server{
		Addr:    fmt.Sprintf("localhost%v", httpPort),
		Handler: mux,
	}
	httpL, err := net.Listen("tcp",
		fmt.Sprintf("localhost%v", httpPort),
	)
	if err != nil {
		log.Fatal(err)
	}
	go func(wg *sync.WaitGroup, listener net.Listener) {
		defer wg.Done()
		if err := serverHTTP.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}(&wg, httpL)
	wg.Wait()
}
