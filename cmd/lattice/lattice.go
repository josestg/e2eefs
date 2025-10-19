package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("PONG!"))
		if err != nil {
			log.Printf("cannot reply: %s", err.Error())
		}
	})

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
