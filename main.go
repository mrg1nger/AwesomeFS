package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Command-line flag for port
	port := flag.Int("port", 8081, "port to host the server on")
	flag.Parse()

	// Serve files from the "./files" directory
	fs := http.FileServer(http.Dir("./files"))
	http.Handle("/", fs)

	// Start the server
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting file server at http://localhost%s\n", addr)
	fmt.Println("Serving files from ./files directory")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
