package main

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed index.html
var index []byte

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", home)
	server := http.Server{Addr: ":8000", Handler: router}
	log.Println("Started server on :8000")
	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(index)
}
