package product

import (
	"context"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpcgateway/internal/pb/product"
	"log"
)

type Server struct {
	product.UnimplementedProductInfoServer
	productMap map[string]*product.Product
}

func (s *Server) AddProduct(ctx context.Context, in *product.Product) (*product.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*product.Product)
	}
	s.productMap[in.Id] = in
	log.Printf("Product %v : %v - Added.", in.Id, in.Name)
	return &product.ProductID{ProductId: in.Id}, status.New(codes.OK, "").Err()
}

func (s *Server) GetProduct(ctx context.Context, in *product.ProductID) (*product.Product, error) {
	log.Printf("Product %v - Get.", in.ProductId)
	product, exists := s.productMap[in.ProductId]
	if exists && product != nil {
		log.Printf("Product %v : %v - Retrieved.", product.Id, product.Name)
		return product, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.ProductId)
}
