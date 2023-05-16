package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
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
	flag.Parse()

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	service, err := service.New(*resourceRoot, *stacksDir)
	if err != nil {
		log.Fatalf("Failed to create LogViz service: %s", err)
	}

	mux := http.DefaultServeMux
	service.RegisterHandlers(mux)
	mux.Handle("/", http.FileServer(http.Dir(*resourceRoot)))
	fmt.Printf("Serving on port %d\n", *port)
	http.ListenAndServe(
		fmt.Sprintf(":%d", *port),
		mux,
	)
}

func CreateCollection(ctx context.Context, client *ent.Client) (*ent.Collection, error) {
	stacks, err := os.ReadFile("cockroachdb_example_snapshot.txt")
	if err != nil {
		return nil, fmt.Errorf("failed reading file: %w", err)
	}

	var snaps []*ent.ProcessSnapshot
	for i := 1; i <= 2; i++ {
		s, err := client.ProcessSnapshot.Create().
			SetID(int64(i)).
			SetProcessID("node-1").
			SetSnapshot(string(stacks)).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating snapshot: %w", err)
		}
		snaps = append(snaps, s)
	}

	c, err := client.Collection.Create().
		SetID(1).
		SetName("crdb-20230516-155400").
		AddProcessSnapshots().
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
