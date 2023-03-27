package main

import (
	"log"
	"net/http"

	"github.com/aajdinov/inventoryservice/database"
	"github.com/aajdinov/inventoryservice/product"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api/v1"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	err := http.ListenAndServe(":5001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
