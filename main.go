// Created by mvp on 14/02/2025

package main

import (
	"fmt"
	"github.com/plitto007/go-cache-wrapper/caching"
	cachingutil "github.com/plitto007/go-cache-wrapper/caching/util"
	"reflect"
)

type TestStruct struct {
	Name string
	Age  int
}

func (t *TestStruct) GetData() string {
	return fmt.Sprintf("Name: %v - Age: %v", t.Name, t.Age)
}

func GetData() string {
	return "Hello"
}

func main() {
	var t TestStruct
	caching.InitCaching(nil)
	var runtimeFuncName = cachingutil.GetRuntimeFuncName(t.GetData)
	fmt.Println("runtimeFuncName:", runtimeFuncName)

	runtimeFuncName = cachingutil.GetRuntimeFuncName(GetData)
	fmt.Println("runtimeFuncName:", runtimeFuncName)

	cacheCompute := cachingutil.FuncCacheDecorator(caching.GetCacheInstance(), GetData, 10)
	fmt.Println("Result: ", reflect.TypeOf(cacheCompute()))

	fmt.Println("Data: ", caching.TriggerFunc(t.GetData))
}
