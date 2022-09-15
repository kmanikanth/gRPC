package intercept

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

func AuthStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	token := md.Get("authorization")
	if len(token) == 0 {
		return status.Errorf(codes.Unauthenticated, "missing authorization token")
	}
	tokens := map[string]bool{"token1": true, "token2": true, "token3": true, "token4": true, "token5": true}
	if _, ok := tokens[token[0]]; !ok {
		log.Printf("method: %v - invalid token: %v", info.FullMethod, token[0])
		return status.Errorf(codes.Unauthenticated, fmt.Sprintf("invalid token: %v", token[0]))
	}
	return handler(srv, stream)
}
