package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aajdinov/inventoryservice/cors"
	"golang.org/x/net/websocket"
)

const productsBasePath = "products"

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	reportHandler := http.HandlerFunc(handleProductReport)
	http.Handle("/websocket", websocket.Handler(productSocket))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))
	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, productsBasePath), cors.Middleware(reportHandler))
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		addProduct(w, r)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/products/"):]
	switch r.Method {
	case http.MethodGet:
		findProductByID(w, r, id)
	case http.MethodPut, http.MethodPatch:
		putProduct(w, r, id)
	case http.MethodDelete:
		deleteProduct(w, r, id)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}
}

func findProductByID(w http.ResponseWriter, r *http.Request, id string) {
	encoder := json.NewEncoder(w)
	productID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := getProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encoder.Encode(product)
}

func putProduct(w http.ResponseWriter, r *http.Request, id string) {
	var product Product
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not update product"))
		return
	}
	if id == strconv.Itoa(0) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.ProductID = productID
	updateProduct(product)

	w.WriteHeader(http.StatusOK)
}

func deleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	productID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	removeProduct(productID)

	w.WriteHeader(http.StatusAccepted)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not add product"))
		return
	}
	if product.ProductID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productID, err := insertProduct(product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(productID)))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	productList, err := getProductList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder.Encode(productList)
}
