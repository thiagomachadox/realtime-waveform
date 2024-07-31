package main

import (
	"log"
	"net/http"

	"github.com/thiagomachadox/realtime-waveform/web"
)

func main() {
	server := web.NewServer()
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
