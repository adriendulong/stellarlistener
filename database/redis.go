package database

import (
	"errors"

	"github.com/mediocregopher/radix.v2/pool"
)

//Redis is a struct that hold the pool connection
type Redis struct {
	Pool *pool.Pool
}

// New open the pool connection and return
func New() (r Redis) {
	p, err := pool.New("tcp", "localhost:6379", 10)
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
