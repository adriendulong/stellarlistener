package database

import (
	"errors"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

//Redis is a struct that hold the pool connection
type Redis struct {
	Client *redis.Client
}

// New open the pool connection and return
func New() (r Redis) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	r = Redis{Client: client}
	return
}

// SetKey simple function that set a key with a value
func (r *Redis) SetKey(k string, v interface{}) {
	if r.Client == nil {
		panic(errors.New("No CLient found"))
	}

	resp := r.Client.Set(k, v, 0)
	if resp.Err() != nil {
		panic(resp.Err())
	}
}
