package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	scheme = "http"
)

func main() {
	addr := flag.String("addr", ":4002", "HTTP Port Address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", pageHandler)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/angelcastaneda.asc", pgpHandler)
	mux.HandleFunc("/pgp", pgpHandler)
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/cv", cvHandler)
	mux.HandleFunc("/claim-before.pdf", cbtsHandler)

	log.Println("Starting server on port ", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
