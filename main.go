package main

import (
	"log"
	"net/http"

	"github.com/aajdinov/inventoryservice/database"
	"github.com/aajdinov/inventoryservice/product"
	"github.com/aajdinov/inventoryservice/receipt"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api/v1"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":5001", nil))
}
