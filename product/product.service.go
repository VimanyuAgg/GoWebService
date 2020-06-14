package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/goWebServices/one-o-one/cors"
)

const productsBasePath = "products"

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsJSON, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJSON)
		return
	case http.MethodPost:
		product := Product{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error happened while reading request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &product)
		if err != nil {
			log.Println("Error happened while unmarshalling request into product")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if product.ProductID != 0 {
			log.Println("The request was setting the ProductID which needs to be set automatically")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = insertNewProduct(product)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments) - 1])
	if err != nil {
		log.Println("Invalid ProductID found in the request")
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

	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(product)
		if err != nil {
			log.Println("Error while converting DB struct to JSON")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(productJSON)

	case http.MethodPut:
		updatedProduct := Product{}
		incomingProduct, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error happened while reading incoming request body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(incomingProduct, &updatedProduct)
		if err != nil {
			log.Println("Error happened while converting incoming request to product struct")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//incoming request body contains the product ID. We need to match if this matches the URL path
		if updatedProduct.ProductID != 0 && updatedProduct.ProductID != productID{
			log.Println("ProductID in request body doesn't match the URL path")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateProduct(updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		err = removeProduct(productID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
