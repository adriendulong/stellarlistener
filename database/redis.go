package database

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/mediocregopher/radix.v2/pool"
	log "github.com/sirupsen/logrus"
)

//Redis is a struct that hold the pool connection
type Redis struct {
	Pool *pool.Pool
}

// New open the pool connection and return
func New() (r Redis) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Info(os.Getenv("REDIS_URL"))
	p, err := pool.New("tcp", os.Getenv("REDIS_URL"), 20)
	if err != nil {
		panic(err)
	}
	r = Redis{Pool: p}
	return
}

// SetKey simple function that set a key with a value
func (r *Redis) SetKey(k string, v interface{}) {
	if r.Pool == nil {
		panic(errors.New("No pool found"))
	}

	resp := r.Pool.Cmd("set", k, v)
	if resp.Err != nil {
		panic(resp.Err)
	}
}
