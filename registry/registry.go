package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"serviceRegistry/services"
)

func main() {
	registry := services.NewRegistry()
	port := os.Getenv("PORT")
	playOnRegistry(port, registry)
	select {}
}

func playOnRegistry(port string, registry *services.Registry) {
	server := rpc.NewServer()
	err := server.RegisterName("Registry", registry)
	if err != nil {
		fmt.Printf("format of Register service is not correct: %s", err)
	}
	port = fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", port)
	go server.Accept(lis)
}
