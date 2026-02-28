package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Command-line flags
	host := flag.String("host", "0.0.0.0", "interface to bind to")
	port := flag.Int("port", 8081, "port to host the server on")
	flag.Parse()

	// Serve files from the "./files" directory
	fs := http.FileServer(http.Dir("./files"))
	http.Handle("/", fs)

	// Start the server
	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Starting file server at http://%s\n", addr)
	fmt.Println("Serving files from ./files directory")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
