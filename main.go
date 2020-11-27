//go:generate go run github.com/99designs/gqlgen generate
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/scayle/common-go"
	"github.com/scayle/gateway/graph/generated"
	"github.com/scayle/gateway/graph/resolver"
	pb "github.com/scayle/proto-go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const defaultPort = "8080"

func main() {
	service := common.GetRandomServiceWithConsul("user-service")
	if service == nil {
		panic("no user-service found")
	}

	host := net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))
	// Set up a connection to the server.
	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Panic the service if the connection gets lost which may trigger a restart by docker compose.
	// TODO: add better error handling in this case as this is the gateway which should not restart if any dependent
	//       service get's lost. Instead it should try to reconnect.
	go func() {
		// continue checking for state change
		// until one of break states is found
		for {
			change := conn.WaitForStateChange(context.Background(), conn.GetState())
			if !change {
				// ctx is done
				// something upstream is cancelling
				panic("grpc lost connection")
			}

			currentState := conn.GetState()
			if currentState == connectivity.Shutdown || currentState == connectivity.TransientFailure {
				panic("grpc lost connection")
			}
		}
	}()

	userService := pb.NewUserServiceClient(conn)

	// ToDo: register also with consul
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			UserService: userService,
		},
		Directives: resolver.Directives(),
	}))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: make configurable
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(Authenticator(userService))

	r.Handle("/", playground.Handler("GraphQL playground", "/v1"))
	r.Handle("/v1", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+defaultPort, r))
}
