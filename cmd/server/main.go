package main

import (
	"log"
	"net"

	"github.com/Likebutter/todo-api/internal/server"
)

func main() {
	port := ":8080"
	log.Printf("Starting listening on port %s", port)

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen to tcp connections: %v", err)
	}

	log.Printf("Listening on port: %s", port)

	server := server.NewGRPCServer()

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
