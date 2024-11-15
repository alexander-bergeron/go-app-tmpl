package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	userpb "github.com/alexander-bergeron/go-app-tmpl/gen/go/proto/user/v1"
	users "github.com/alexander-bergeron/go-app-tmpl/internal/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("error running application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("closing server gracefully")
}

func run(ctx context.Context) error {
	// initialize empty serverOptions
	var serverOpts []grpc.ServerOption

	// specify grpc server port
	const addr = ":9090"

	// set conn string
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
	}
	defer db.Close()

	// add tls
	tlsCredentials, err := credentials.NewServerTLSFromFile("/certs/server.crt", "/certs/server.key")
	if err != nil {
		return fmt.Errorf("failed to load tls credentials: %w", err)
	}
	serverOpts = append(serverOpts, grpc.Creds(tlsCredentials))

	// initialize server
	s := grpc.NewServer(serverOpts...)
	userpb.RegisterUserServiceServer(s, users.NewUserService(db))

	// Enable reflection
	reflection.Register(s)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on address %q: %w", addr, err)
		}

		slog.Info("starting grpc server on address", slog.String("address", addr))

		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve grpc service: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		// Start HTTP gateway
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()

		// credentials for grpc gateway (client) to talk to grpc service
		// this verifies the servers identity with the ca.crt
		clientCreds, err := credentials.NewClientTLSFromFile(
			"/certs/ca.crt",
			"localhost",
		)
		if err != nil {
			return fmt.Errorf("failed to create client credentials: %w", err)
		}

		opts := []grpc.DialOption{grpc.WithTransportCredentials(clientCreds)}
		// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		if err := userpb.RegisterUserServiceHandlerFromEndpoint(
			ctx,
			mux,
			"localhost:9090",
			opts,
		); err != nil {
			log.Fatalf("failed to register gateway: %v", err)
		}

		handler := corsMiddleware(mux)

		// Serve over tls
		slog.Info("Starting HTTPS server on :8080")
		if err := http.ListenAndServeTLS(
			":8080",
			"/certs/server.crt",
			"/certs/server.key",
			handler,
		); err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}

		// log.Printf("Starting HTTP server on :8080")
		// if err := http.ListenAndServe(":8080", handler); err != nil {
		// 	log.Fatalf("failed to serve: %v", err)
		// }
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		s.GracefulStop()

		return nil
	})

	return g.Wait()
}
