package main

import (
	"context"
	"entgo.io/contrib/entgql"
	"flag"
	"fmt"
	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"path"
	server "stacksviz"
	"stacksviz/datasource"
	"stacksviz/ent"
	"stacksviz/util"
)

var (
	port         = flag.Int("port", 7410, "Port to serve LogViz clients on")
	resourceRoot = flag.String("resource_root", "", "The path to the LogViz tool client resources")
	stacksDir    = flag.String("stacks_dir", ".", "The root path for visualizable stacks")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	conf, err := util.ReadConfig(path.Join(cwd, "config.yaml"))
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	// !!! client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client, err := ent.Open("sqlite3", "file:stackviz.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	client = client.Debug()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	stacksFetcher := datasource.NewStacksFetcher(client)
	// !!!
	//service, err := service.New(*resourceRoot, stacksFetcher)
	//if err != nil {
	//	log.Fatalf("Failed to create LogViz service: %s", err)
	//}

	mux := http.DefaultServeMux
	// !!! service.RegisterHandlers(mux)
	mux.Handle("/", http.FileServer(http.Dir(*resourceRoot)))

	// Create the Graphql server and register it and the playground.
	graphqlServer := graphqlhandler.NewDefaultServer(server.NewSchema(client, stacksFetcher, conf))
	graphqlServer.Use(entgql.Transactioner{TxOpener: client})
	mux.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	mux.Handle("/graphql", graphqlServer)

	// Start the HTTP server.
	fmt.Printf("Serving on port %d. Go to http://localhost:7410 for the app "+
		"and http://localhost:7410/playground for a GraphQL playground.\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
}
