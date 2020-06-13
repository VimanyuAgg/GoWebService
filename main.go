package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type fooHandler struct {
	Message string
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(f.Message))
}

func barHandler(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Bar Called"))
}

type Product struct {
	ProductID int `json."productId"`
	Manufacturer string `json."manufacturer"`
	Sku string `json."sku"`
	Upc string `json."upc"`
	PricePerUnit string `json."pricePerUnit"`
	QuantityAtHand int `json."quantityAtHand"`
	ProductName string `json."productName"`
}

var productList []Product

func init() {
	productsJSON := `[
		{
			"productId": 1,
			"manufacturer": "Charizard Company Inc.",
			"pricePerUnit": "10.01",
			"sku": "CHARZ8000",
			"upc": "CH3214",
			"quantityAtHand": 2434,
			"productName": "Swirling Inferno"
		},
		{
			"productId": 2,
			"manufacturer": "Bulba Association Ltd.",
			"pricePerUnit": "2.40",
			"sku": "BULBA01Z",
			"upc": "BLB11A",
			"quantityAtHand": 50,
			"productName": "Big Hose"
		},
		{
			"productId": 3,
			"manufacturer": "Pikachu Product Ltd",
			"pricePerUnit": "50",
			"sku": "PIKAPIKA-1",
			"upc": "PIKA01A",
			"quantityAtHand": 10,
			"productName": "Blazing Bolt 5000"
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
			log.Println("Error happened while unmarshalling request into Product")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if product.ProductID != 0 {
			log.Println("The request was setting the ProductID which needs to be set automatically")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product.ProductID = getNextID()
		productList = append(productList, product)

		w.WriteHeader(http.StatusCreated)
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
	product, listItemIndex := findProductByID(productID)

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
			log.Println("Error happened while converting incoming request to Product struct")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//incoming request body contains the Product ID. We need to match if this matches the URL path
		if updatedProduct.ProductID != 0 && updatedProduct.ProductID != productID{
			log.Println("ProductID in request body doesn't match the URL path")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updatedProduct.ProductID = productID
		productList[listItemIndex] = updatedProduct
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler, middleware start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("after handler, middleware end - Time taken: %s \n", time.Since(start))
	})
}

func main() {
	fmt.Println("Running main.go...")

	//Uncomment Below for handling requests without middleware

	//http.Handle("/foo", &fooHandler{"Foo called!"})
	//http.HandleFunc("/bar", barHandler)
	//http.HandleFunc("/products", productsHandler)
	//http.HandleFunc("/products/",productHandler)

	//With Middleware
	http.Handle("/foo", middlewareHandler(&fooHandler{"Foo called!"}));

	barHandler1 := http.HandlerFunc(barHandler)
	http.Handle("/bar", middlewareHandler(barHandler1))

	productItemHandler := http.HandlerFunc(productHandler)
	productListHandler := http.HandlerFunc(productsHandler)

	http.Handle("/products", middlewareHandler(productListHandler))
	http.Handle("/products/", middlewareHandler(productItemHandler))

	http.ListenAndServe(":5000", nil)
}

func getNextID() int{
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}
	return highestID + 1
}

func findProductByID(productID int) (*Product, int) {
	for i, product := range productList {
		if product.ProductID == productID {
			return &product, i
		}
	}
	return nil, 0
}
