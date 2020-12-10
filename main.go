package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const addr = ":8080"
	fmt.Printf("Start Listen %s", addr)

	http.HandleFunc("/", helloWorldHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Ironhead")
}
