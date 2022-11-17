package main

import (
	"fmt"
	"log"
	"net/http"
)

func NapHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Nap time!")
}

func main() {
	http.HandleFunc("/", NapHandler)
	log.Println("Napping...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
