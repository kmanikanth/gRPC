package product

import (
	"bidirectional-streaming/internal/pb/product"
	"context"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
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
	return &product.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

func (s *Server) GetProduct(ctx context.Context, in *product.ProductID) (*product.Product, error) {
	product, exists := s.productMap[in.Value]
	if exists && product != nil {
		log.Printf("Product %v : %v - Retrieved.", product.Id, product.Name)
		return product, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}

func (s *Server) GetProducts(req *product.ProductsRequest, stream product.ProductInfo_GetProductsServer) error {
	for _, prd := range s.productMap {
		err := stream.Send(prd)
		if err != nil {
			log.Printf("error sending message to stream : %v", err)
			return err
		}
		log.Printf("Product %v : %v - sent.", prd.Id, prd.Type)
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

func (s *Server) BulkAdd(stream product.ProductInfo_BulkAddServer) error {
	var wg sync.WaitGroup
	ch := make(chan *product.BulkAddResponse)
	wg.Add(2)
	go s.consumer(&wg, ch, stream)
	go s.producer(&wg, ch, stream)
	wg.Wait()
	return nil
}

func (s *Server) consumer(wg *sync.WaitGroup, ch chan<- *product.BulkAddResponse, stream product.ProductInfo_BulkAddServer) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		out, _ := uuid.NewV4()
		res.Id = out.String()
		if s.productMap == nil {
			s.productMap = make(map[string]*product.Product)
		}
		s.productMap[res.Id] = res
		log.Printf("Product %v : %v - Added.", res.Id, res.Name)
		ch <- &product.BulkAddResponse{Id: res.Id, Name: res.Name}
	}
	close(ch)
}

func (s *Server) producer(wg *sync.WaitGroup, ch <-chan *product.BulkAddResponse, stream product.ProductInfo_BulkAddServer) {
	defer wg.Done()
	for res := range ch {
		err := stream.Send(res)
		if err != nil {
			log.Fatalf("Error while sending stream: %v", err)
		}
	}
}
