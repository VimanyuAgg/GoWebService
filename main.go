package main

import (
	"fmt"
	"github.com/goWebServices/one-o-one/database"
	"log"
	"net/http"
	"time"

	"github.com/goWebServices/one-o-one/product"
)

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler, middleware start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("after handler, middleware end - Time taken: %s \n", time.Since(start))
	})
}

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	err := http.ListenAndServe(":5000", nil)

	if err != nil {
		log.Fatal(err)
	}
}
