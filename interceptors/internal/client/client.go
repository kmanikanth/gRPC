package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"interceptors/internal/pb/product"
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

func (c *Client) AddProduct() {
	for i := 1; i <= 10; i++ {
		ctx := context.Background()
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("token%v", i))
		productID, err := c.client.AddProduct(ctx, &product.Product{Name: fmt.Sprintf("Product %v", i)})
		if err != nil {
			log.Printf("Error while calling AddProduct RPC: %v", err)
			continue
		}
		log.Printf("Product ID: %v", productID.Value)
	}
}

func (c *Client) GetProducts() {
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("token%v", 8))
	stream, err := c.client.GetProducts(ctx, &product.ProductsRequest{})
	if err != nil {
		log.Printf("Error while calling GetProducts RPC: %v", err)
		return
	}
	log.Printf("request sent to get products")
	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error while receiving stream: %v", err)
			//st, _ := status.FromError(err)
			//log.Printf("Error while receiving stream status: %v", st.Err())
			break
		}
		log.Printf("Product: %v : %v - received", product.GetId(), product.GetName())
	}
}
