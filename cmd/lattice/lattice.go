package main

import (
	"log"
	"net/http"
	"os"
)

// Adapter Pattern
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("PONG!"))
		if err != nil {
			log.Printf("cannot reply: %s", err.Error())
		}
	})

	f := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("PONG!"))
		if err != nil {
			log.Printf("cannot reply: %s", err.Error())
		}
	}

	mux.Handle("/echo", HandlerFunc(f))

	srv := http.Server{
		Addr:    "localhost:8080", // host:port
		Handler: mux,
	}

	log.Printf("server is listening: %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println("error:", err.Error())
		os.Exit(1)
	}
}
