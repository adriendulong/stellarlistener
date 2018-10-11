package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	c "github.com/adriendulong/go/stellar/controller"
	"github.com/adriendulong/go/stellar/database"
	"github.com/gorilla/mux"
)

var redis *database.Redis

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// HomeHandler is a simple Handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	nbOperations, err := c.TotalOperationsOfDay(time.Now(), redis)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	sResp := fmt.Sprintf("Number total of operation on %s is %d", time.Now(), nbOperations)
	fmt.Fprintln(w, sResp)
}

func main() {
	r := database.New()
	redis = &r

	fmt.Println("Coucou")

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":8080", router))
}
