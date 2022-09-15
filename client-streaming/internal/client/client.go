package client

import (
	"client-streaming/internal/pb/product"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
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
		p := &product.Product{
			Name:        fmt.Sprintf("Product %v", i),
			Description: "XYZ",
			Type:        product.ProductType_PRODUCT_TYPE_FURNITURE,
		}
		productID, err := c.client.AddProduct(context.Background(), p)
		if err != nil {
			log.Fatalf("Error while calling AddProduct RPC: %v", err)
		}
		p.Id = productID.GetValue()
		c.productMap[productID.GetValue()] = p
	}
}

func (c *Client) GetProducts() {
	stream, err := c.client.GetProducts(context.Background(), &product.ProductsRequest{})
	if err != nil {
		log.Fatalf("Error while calling GetProducts RPC: %v", err)
	}
	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving stream: %v", err)
		}
		log.Printf("Product Type: %v", product.GetType())
		// do some bussiness operation
		//time.Sleep(10 * time.Second)
	}
}

func (c *Client) UpdateProducts() {
	stream, err := c.client.UpdateProducts(context.Background())
	if err != nil {
		log.Fatalf("Error while calling UpdateProducts RPC: %v", err)
		return
	}
	for _, product := range c.productMap {
		product.Description = product.Description + " - Updated Description"
		err := stream.Send(product)
		if err != nil {
			log.Fatalf("Error while sending update: %v", err)
			return
		}
		log.Printf("Requested Update for Product ID: %v", product.Id)

	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
		return
	}
	log.Printf("UpdateProducts Response: %v", resp)
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
