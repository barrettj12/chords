package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/barrettj12/chords/gqlgen"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: &gqlgen.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("error from net.Listen: %v", err)
	}
	log.Printf("GraphQL playground serving on %v", ln.Addr())
	log.Fatal(http.Serve(ln, nil))
}
