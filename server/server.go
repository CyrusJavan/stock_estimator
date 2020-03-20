package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CyrusJavan/stock_estimator/simulation"
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

	r.HandleFunc("/invest/{stock}/", investHandler).
		Queries(
			"principal", "{principal:[0-9]+}",
			"recurring", "{recurring:[0-9]+}",
			"start", "{start:[0-9]{4}-[0-9]{2}-[0-9]{2}}",
			"end", "{end:[0-9]{4}-[0-9]{2}-[0-9]{2}}")

	http.Handle("/", r)
	addr := "0.0.0.0"
	port := "8000"
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", addr, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server started on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func investHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !validStock(vars["stock"]) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid stock name: %s", vars["stock"])))
		return
	}
	sim, err := simulation.NewSimulation(stockNameToDataFile(vars["stock"]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	principal, err := strconv.ParseFloat(vars["principal"], 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	recurringInvestment, err := strconv.ParseFloat(vars["recurring"], 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	worth, invested := sim.InvestOverTime(
		vars["start"],
		vars["end"],
		principal,
		recurringInvestment,
	)

	type investResponse struct {
		Worth    float64 `json:"worth"`
		Invested float64 `json:"invested"`
	}

	resp, err := json.Marshal(investResponse{
		worth,
		invested,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func stockNameToDataFile(name string) string {
	nameToFile := map[string]string{
		"djia": "data/DJI.csv",
	}

	return nameToFile[name]
}

func validStock(name string) bool {
	if name != "djia" {
		return false
	}
	return true
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
