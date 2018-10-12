package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	c "github.com/adriendulong/go/stellar/controller"
	"github.com/adriendulong/go/stellar/database"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
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

	//Get total operations of the day
	nbOperations, err := c.TotalOperationsOfDay(time.Now(), redis)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	types := make(map[string](chan int))
	types["payment"] = make(chan int)
	types["manage_offer"] = make(chan int)
	types["inflation"] = make(chan int)

	nbOpeTypes := make(map[string]int)
	nbOpeTypes["payment"] = 0
	nbOpeTypes["manage_offer"] = 0
	nbOpeTypes["inflation"] = 0

	for k, v := range types {
		go c.TotalOpearationsOfDayType(time.Now(), k, redis, v)
	}

	for i := 0; i < len(nbOpeTypes); i++ {
		select {
		case nbOpeTypes["payment"] = <-types["payment"]:
			log.Info("Got payment nb")
		case nbOpeTypes["manage_offer"] = <-types["manage_offer"]:
			log.Info("Got manage_offer nb")
		case nbOpeTypes["inflation"] = <-types["inflation"]:
			log.Info("Got inflation nb")
		}
	}

	w.WriteHeader(http.StatusOK)
	sResp := fmt.Sprintf("Number total of operation on %s %dth is %d\nNumber of Payment is %d\nNumber of Manage Offer is %d\nNumber of Inlfation is %d", time.Now().Month().String(), time.Now().Day(), nbOperations, nbOpeTypes["payment"], nbOpeTypes["manage_offer"], nbOpeTypes["inflation"])
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
