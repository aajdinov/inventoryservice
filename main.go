package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func init() {
	productsJSON := `[
		{
			"productId":      1,
			"manufacturer":   "Johns-Jenkins",
			"sku":            "p5z343vdS",
			"upc":            "939581000000",
			"pricePerUnit":   "497.45",
			"quantityOnHand": 9703,
			"productName":    "sticky note"
		},
		{
			"productId":      2,
			"manufacturer":   "Hessel, Schimmel and Feeney",
			"sku":            "i7v300kmx",
			"upc":            "740979000000",
			"pricePerUnit":   "282.29",
			"quantityOnHand": 9217,
			"productName":    "leg warmers"
		},
		{
			"productId":      3,
			"manufacturer":   "Swaniawski, Bartoletti and Bruen",
			"sku":            "q0L657ys7",
			"upc":            "111730000000",
			"pricePerUnit":   "436.26",
			"quantityOnHand": 5905,
			"productName":    "lamp shade"
		},
		{
			"productId":      4,
			"manufacturer":   "Big Box Company",
			"sku":            "456qHJK",
			"upc":            "414654444566",
			"pricePerUnit":   "$12.99",
			"quantityOnHand": 28,
			"productName":    "Gizmo"
		}
	]`
	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		addProduct(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/products/"):]
	switch r.Method {
	case http.MethodGet:
		findProductByID(w, r, id)
	case http.MethodPut, http.MethodPatch:
		updateProduct(w, r, id)
	case http.MethodDelete:
		deleteProduct(w, r, id)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

func findProductByID(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	for _, product := range productList {
		if strconv.Itoa(product.ProductID) == id {
			encoder.Encode(product)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func updateProduct(w http.ResponseWriter, r *http.Request, id string) {
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

	for i, p := range productList {
		if strconv.Itoa(p.ProductID) == id {
			id_, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			product.ProductID = id_
			productList[i] = product
			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			encoder.Encode(productList[i])
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	for i, p := range productList {
		if strconv.Itoa(p.ProductID) == id {
			productList = append(productList[:i], productList[i+1:]...)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func getNextID() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}

	return highestID + 1
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

	product.ProductID = getNextID()
	productList = append(productList, product)
	w.WriteHeader(http.StatusCreated)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(productList)
}

func main() {
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)

	http.ListenAndServe(":5000", nil)
}
