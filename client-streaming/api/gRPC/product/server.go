package product

import (
	"client-streaming/internal/pb/product"
	"context"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
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
	return &product.ProductID{Value: in.Id}, nil
}

func (s *Server) GetProduct(ctx context.Context, in *product.ProductID) (*product.Product, error) {
	product, exists := s.productMap[in.Value]
	if exists && product != nil {
		return product, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}

func (s *Server) GetProducts(req *product.ProductsRequest, stream product.ProductInfo_GetProductsServer) error {
	for _, prd := range s.productMap {
		err := stream.Send(prd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) UpdateProducts(stream product.ProductInfo_UpdateProductsServer) error {
	var prds string
	for {
		prd, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&product.UpdateResponse{Message: "Products Updated Successfully: " + prds})
		}
		if err != nil {
			log.Printf("error sending message to stream : %v", err)
			return err
		}
		if _, ok := s.productMap[prd.GetId()]; ok {
			s.productMap[prd.GetId()] = prd
			prds += prd.GetId() + ", "
			log.Printf("Product %v : %v - Updated.", prd.Id, prd.Name)
		}
	}
}
