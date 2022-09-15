package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	productRPC "server-streaming/api/gRPC/product"
	"server-streaming/internal/pb/product"
)

const (
	port = ":8080"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	product.RegisterProductInfoServer(s, &productRPC.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
