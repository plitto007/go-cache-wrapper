// Created by mvp on 14/02/2025

package util

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func GetRuntimeFuncName(i interface{}) string {
	fmt.Println("Get runtime function name")
	fn := reflect.ValueOf(i)
	fmt.Printf("Function type %v\n", fn.Kind())
	if fn.Kind() == reflect.Ptr {
		fn = fn.Elem()
	}
	funcPtr := fn.Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	return funcName
}

// GenerateFuncCacheKey will generate the cache key for the function name (which created by GetRuntimeFuncName)
// Format of cache key: functionName-->args(separated by `|`)
func GenerateFuncCacheKey(funcName string, args []reflect.Value) string {
	var sb strings.Builder
	sb.WriteString(funcName + "-->")
	for _, arg := range args {
		sb.WriteString(fmt.Sprintf("%v|", arg.Interface()))
	}
	return sb.String()
}

// FuncCacheDecorator create Generic function caching decorator
func FuncCacheDecorator[T any](cache *ristretto.Cache, fn T, ttl time.Duration) func(args ...interface{}) interface{} {
	fnValue := reflect.ValueOf(fn)
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		panic("FuncCacheDecorator expects a function")
	}

	fnName := GetRuntimeFuncName(fn)

	return func(args ...interface{}) interface{} {
		// converts args to reflect.Value
		var reflectArgs []reflect.Value
		for _, arg := range args {
			reflectArgs = append(reflectArgs, reflect.ValueOf(arg))
		}

		// Generate cache key
		key := GenerateFuncCacheKey(fnName, reflectArgs)
		fmt.Printf("Cache key:%v\n", key)
		// Check cache key
		if val, found := cache.Get(key); found {
			return val.([]reflect.Value)
		}
		// Call the origin function
		results := fnValue.Call(reflectArgs)

		// Store  result in cache with TTL
		cache.SetWithTTL(key, results, 1, ttl)
		cache.Wait()
		return results[0].Interface()
	}
	//// wrap the function
	//wrappedFunc := reflect.MakeFunc(fnType, func(args []reflect.Value) (results []reflect.Value) {
	//	key := GenerateFuncCacheKey(fnName, args)
	//
	//	//var cache = caching.GetCacheInstance()
	//	// Check cache key
	//	if val, found := cache.Get(key); found {
	//		return val.([]reflect.Value)
	//	}
	//
	//	// Call the origin function
	//	results = fnValue.Call(args)
	//
	//	// Store  result in cache with TTL
	//	cache.SetWithTTL(key, results, 1, ttl)
	//	cache.Wait()
	//	return results
	//})
	//return wrappedFunc.Interface().(T)
}
