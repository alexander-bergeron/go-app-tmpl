package main

import (
	"log"
	"net/http"

	handler "github.com/alexander-bergeron/go-app-tmpl/internal/frontend"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log.Println("Starting Webserver")

	// Create a new ServeMux
	mux := http.NewServeMux()

	// credentials for grpc gateway (client) to talk to grpc service
	// this verifies the servers identity with the ca.crt
	clientCreds, err := credentials.NewClientTLSFromFile(
		"/certs/ca.crt",
		"localhost",
	)
	if err != nil {
		log.Fatal("failed to create client credentials: %w", err)
	}

	conn, err := grpc.NewClient("proto:9090",
		grpc.WithTransportCredentials(clientCreds),
	)
	if err != nil {
		log.Fatal("Failed to connect to gRPC server:", err)
	}
	defer conn.Close()

	// Create a new ListingHandler
	listingHandler, err := handler.NewPageHandler(conn)
	if err != nil {
		log.Fatal("Failed to create ListingHandler:", err)
	}

	listingHandler.RegisterRoutes(mux)

	// Start the HTTP server using the mux
	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
