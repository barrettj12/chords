package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/barrettj12/chords/gqlgen"
	"github.com/barrettj12/chords/src/dblayer"
)

func main() {
	// Try to read logging flags from LOG_FLAGS environment variable
	// Invalid/unset values will just default to 0 (no flags)
	flags, _ := strconv.Atoi(os.Getenv("LOG_FLAGS"))

	// Initialise logger
	logger := log.New(os.Stdout, "", flags)

	// Set up DB
	dbURL := os.Getenv("DATABASE_URL")
	db := dblayer.GetDBv1(dbURL, logger)

	srv := handler.NewDefaultServer(gqlgen.NewExecutableSchema(gqlgen.Config{
		Resolvers: &gqlgen.Resolver{
			DB: db,
		}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("error from net.Listen: %v", err)
	}
	log.Printf("GraphQL playground serving on %v", ln.Addr())
	log.Fatal(http.Serve(ln, nil))
}
