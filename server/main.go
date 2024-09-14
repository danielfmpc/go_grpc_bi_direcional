package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/danielfmpc/client_go_rpc_bi_direcional/server/src/pb/shoppingcart"
	"google.golang.org/grpc"
)

type Server struct {
	shoppingcart.ShoppingCartServiceServer
}

func (s *Server) AddItem(srv shoppingcart.ShoppingCartService_AddItemServer) error {
	var quantityItems int32 = 0
	var priceTotal float64 = 0.0
	for {
		newItem, err := srv.Recv()

		if err == io.EOF {
			return srv.Send(&shoppingcart.ProductResponse{
				QuantityItems: quantityItems,
				PriceTotal:    priceTotal,
			})
		}

		if err != nil {
			return fmt.Errorf("error on recv %v", err)
		}

		quantityItems += newItem.GetQuantity()
		priceTotal += (newItem.GetPriceUnit() * float64(newItem.GetQuantity()))

		if err := srv.Send(&shoppingcart.ProductResponse{
			QuantityItems: quantityItems,
			PriceTotal:    priceTotal,
		}); err != nil {
			return fmt.Errorf("erro on send %v", err)
		}
	}
}

func main() {
	fmt.Println("Start")
	lister, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("erro listener ", err)
	}

	s := grpc.NewServer()
	shoppingcart.RegisterShoppingCartServiceServer(s, &Server{})

	if err := s.Serve(lister); err != nil {
		log.Fatalln("err on server ", err)
	}
}
