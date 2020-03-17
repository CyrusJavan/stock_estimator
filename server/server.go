package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// APIServer is
type APIServer struct{}

// NewAPIServer creates a new APIServer
func NewAPIServer() APIServer {
	return APIServer{}
}

// StartServer starts the server
func (s APIServer) StartServer() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/invest", investHandler)
	http.Handle("/", r)
	port := "8000"
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server started on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func investHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hi!")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
