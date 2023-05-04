package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"stacksviz/service"
)

var (
	port         = flag.Int("port", 7410, "Port to serve LogViz clients on")
	resourceRoot = flag.String("resource_root", "", "The path to the LogViz tool client resources")
	stacksDir    = flag.String("stacks_dir", ".", "The root path for visualizable stacks")
)

func main() {
	flag.Parse()

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
