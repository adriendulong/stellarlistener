package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/adriendulong/go/stellar/database"
	"github.com/adriendulong/go/stellar/router"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	re := database.New()
	redis := &(re)
	defer redis.Client.Close()

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(router.RedisMiddleware(redis))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/operation", router.OperationRoutes())

	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, r)
}
