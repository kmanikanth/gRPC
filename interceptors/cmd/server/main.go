package main

import (
	"google.golang.org/grpc"
	productRPC "interceptors/api/gRPC/product"
	"interceptors/internal/intercept"
	"interceptors/internal/pb/product"
	"log"
	"net"
)

const (
	port = ":8080"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(intercept.AuthUnaryInterceptor), grpc.StreamInterceptor(intercept.AuthStreamInterceptor))
	product.RegisterProductInfoServer(s, &productRPC.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
