package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	c "github.com/adriendulong/go/stellar/controller"
	"github.com/adriendulong/go/stellar/database"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	nbOperations, err := c.TotalOperationsOfDay(time.Now(), redis)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	sResp := fmt.Sprintf("Number total of operation on %s is %d", time.Now(), nbOperations)
	fmt.Fprintln(w, sResp)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	r := database.New()
	redis = &r

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.Use(loggingMiddleware)

	addr := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(addr, router))
}
