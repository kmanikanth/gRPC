package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"server-streaming/internal/pb/product"
)

type Client struct {
	client     product.ProductInfoClient
	productMap map[string]*product.Product
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		client:     product.NewProductInfoClient(conn),
		productMap: map[string]*product.Product{},
	}
}

func (c *Client) AddProducts() {
	for i := 1; i <= 10; i++ {
		productID, err := c.client.AddProduct(context.Background(), &product.Product{
			Name:        fmt.Sprintf("Product %v", i),
			Description: "XYZ",
			Type:        product.ProductType_PRODUCT_TYPE_FURNITURE,
		})
		if err != nil {
			log.Fatalf("Error while calling AddProduct RPC: %v", err)
		}
		c.productMap[productID.Value] = &product.Product{Name: fmt.Sprintf("Product %v", i)}
	}
}

func (c *Client) GetProducts() {
	stream, err := c.client.GetProducts(context.Background(), &product.ProductsRequest{})
	if err != nil {
		log.Fatalf("Error while calling GetProducts RPC: %v", err)
	}
	log.Printf("request sent to get products")
	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving stream: %v", err)
		}
		log.Printf("Product: %v : %v - received", product.GetId(), product.GetName())
	}
}

// multiple workers processing at same time
//
//func (c *Client) UpdateFurnitureProducts() {
//	stream, err := c.client.GetProducts(context.Background(), &product.ProductsRequest{Type: product.ProductType_PRODUCT_TYPE_FURNITURE})
//	if err != nil {
//		log.Fatalf("Error while calling GetProducts RPC: %v", err)
//	}
//	var wg sync.WaitGroup
//	ch := make(chan *product.Product)
//	wg.Add(1)
//	go func(wg *sync.WaitGroup, stream product.ProductInfo_GetProductsClient, ch chan<- *product.Product) {
//		defer wg.Done()
//		for {
//			product, err := stream.Recv()
//			if err == io.EOF {
//				close(ch)
//				break
//			}
//			if err != nil {
//				log.Fatalf("Error while receiving stream: %v", err)
//			}
//			ch <- product
//		}
//	}(&wg, stream, ch)
//
//	for i := 0; i < 5; i++ {
//		wg.Add(1)
//		go func(worker int, wg *sync.WaitGroup, ch <-chan *product.Product) {
//			defer wg.Done()
//			for p := range ch {
//				// do some bussiness operation
//				log.Printf("Worker : %d , Product ID: %v", worker, p.GetId())
//			}
//		}(i, &wg, ch)
//	}
//	wg.Wait()
//}
