package router

import (
	"context"
	"net/http"

	"github.com/adriendulong/go/stellar/database"
)

type key int

const rediskey key = 0

// RedisMiddleware is a middleware that sets the redis address to a redis key
func RedisMiddleware(redis *database.Redis) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(SetRedis(r, redis)))
		}
		return http.HandlerFunc(fn)
	}
}

// GetRedis returns a value for this package from the request values.
func GetRedis(r *http.Request) *database.Redis {
	if rv := r.Context().Value(rediskey); rv != nil {
		return rv.(*database.Redis)
	}
	return nil
}

// SetRedis sets a value for this package in the request values.
func SetRedis(r *http.Request, val *database.Redis) context.Context {
	return context.WithValue(r.Context(), rediskey, val)
}
