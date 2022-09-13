package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	"unary/internal/pb/product"
)

const (
	address = "localhost:8080"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := product.NewProductInfoClient(conn)
	for i := 1; i <= 10; i++ {
		productID, err := c.AddProduct(context.Background(), &product.Product{Name: fmt.Sprintf("Product %v", i)})
		if err != nil {
			log.Fatalf("Error while calling AddProduct RPC: %v", err)
		}
		log.Printf("Product ID: %v", productID.Value)
		time.Sleep(time.Second * 2)
	}
}
