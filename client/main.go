package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/danielfmpc/client_go_rpc_bi_direcional/client/src/pb/shoppingcart"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("erro connect to Server ", err)
	}

	defer conn.Close()

	client := shoppingcart.NewShoppingCartServiceClient(conn)

	stream, err := client.AddItem(context.Background())
	if err != nil {
		log.Fatalln("erro channel to stream ", err)
	}

	watch := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(watch)
				return
			}
			if err != nil {
				log.Fatalln("erro rcv ", err)
			}
			fmt.Printf("response %+v\n", response)
		}
	}()

	items := []shoppingcart.ProductRequest{
		{ProductId: 1, Quantity: 2, PriceUnit: 5.0},
		{ProductId: 2, Quantity: 1, PriceUnit: 10.0},
		{ProductId: 3, Quantity: 3, PriceUnit: 55.0},
		{ProductId: 4, Quantity: 1, PriceUnit: 40.23},
	}

	for _, product := range items {
		if err := stream.Send(&product); err != nil {
			log.Fatalln("erro send ", err)
		}
		fmt.Printf("-> send %+v\n", product)
	}

	stream.CloseSend()

	<-watch
}
