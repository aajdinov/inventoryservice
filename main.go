package main

import (
	"log"
	"net/http"

	"github.com/aajdinov/inventoryservice/product"
)

const apiBasePath = "/api/v1"

func main() {
	product.SetupRoutes(apiBasePath)
	err := http.ListenAndServe(":5001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
