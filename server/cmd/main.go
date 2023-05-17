package main

import (
	"context"
	"flag"
	"fmt"
	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	server "stacksviz"
	"stacksviz/ent"
	"stacksviz/ent/collection"
	"stacksviz/service"
)

var (
	port         = flag.Int("port", 7410, "Port to serve LogViz clients on")
	resourceRoot = flag.String("resource_root", "", "The path to the LogViz tool client resources")
	stacksDir    = flag.String("stacks_dir", ".", "The root path for visualizable stacks")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	if _, err := CreateCollection(ctx, client); err != nil {
		log.Fatal(err)
	}

	service, err := service.New(*resourceRoot, *stacksDir)
	if err != nil {
		log.Fatalf("Failed to create LogViz service: %s", err)
	}

	mux := http.DefaultServeMux
	service.RegisterHandlers(mux)
	mux.Handle("/", http.FileServer(http.Dir(*resourceRoot)))

	// Create the Graphql server and register it and the playground.
	graphqlServer := graphqlhandler.NewDefaultServer(server.NewSchema(client))
	mux.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	mux.Handle("/graphql", graphqlServer)

	// Start the HTTP server.
	fmt.Printf("Serving on port %d. Go to http://localhost:7410 for the app "+
		"and http://localhost:7410/playground for a GraphQL playground.\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
}

func CreateCollection(ctx context.Context, client *ent.Client) (*ent.Collection, error) {
	stacks, err := os.ReadFile("datasource/cockroachdb_example_snapshot.txt")
	if err != nil {
		return nil, fmt.Errorf("failed reading file: %w", err)
	}

	var snaps []*ent.ProcessSnapshot
	for i := 1; i <= 2; i++ {
		s, err := client.ProcessSnapshot.Create().
			SetProcessID("node-1").
			SetSnapshot(string(stacks)).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating snapshot: %w", err)
		}
		snaps = append(snaps, s)
	}

	c, err := client.Collection.Create().
		SetName("crdb-20230516-155400").
		AddProcessSnapshots(snaps...).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating collection: %w", err)
	}
	log.Println("collection was created: ", c)

	return c, nil
}

func QueryCollection(ctx context.Context, client *ent.Client) (*ent.Collection, error) {
	c, err := client.Collection.
		Query().
		Where(collection.ID(1)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("collection returned: ", c)
	return c, nil
}
