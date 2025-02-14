// Created by mvp on 14/02/2025

package caching

import (
	"github.com/dgraph-io/ristretto"
	"github.com/plitto007/go-cache-wrapper/caching/util"
	"time"
)

var cache *ristretto.Cache

// initCaching Init the cache, must be called
func InitCaching(cacheConfig *ristretto.Config) {
	var err error
	if cacheConfig == nil {
		// Init cache with default
		cache, err = ristretto.NewCache(&ristretto.Config{
			NumCounters: 10000,     // Number of keys to track = 1000 keys *10(counter)
			MaxCost:     100000000, // Max cache size = 100MB
			BufferItems: 64,        // Performance tuning
		})
	} else {
		cache, err = ristretto.NewCache(cacheConfig)
	}
	if err != nil {
		panic(err)
	}
}

func GetCacheInstance() *ristretto.Cache {
	return cache
}

func TriggerFunc[T any](fun T, args ...interface{}) interface{} {
	return TriggerFuncWithTTL(fun, 0, args...)
}

func TriggerFuncWithTTL[T any](fun T, ttl time.Duration, args ...interface{}) interface{} {
	cacheCompute := util.FuncCacheDecorator(cache, fun, ttl)
	return cacheCompute(args...)
}
