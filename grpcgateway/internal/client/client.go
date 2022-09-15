package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcgateway/internal/pb/product"
	"log"
	"time"
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

func (c *Client) AddProduct() {
	for i := 1; i <= 10; i++ {
		productID, err := c.client.AddProduct(context.Background(), &product.Product{Name: fmt.Sprintf("Product %v", i)})
		if err != nil {
			log.Fatalf("Error while calling AddProduct RPC: %v", err)
		}
		log.Printf("Product ID: %v", productID.ProductId)
		time.Sleep(time.Second * 2)
	}
}
