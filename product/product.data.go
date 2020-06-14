package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int] Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d Products Loaded...\n", len(productMap.m))
}

func loadProductMap() (map[int] Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist ", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)

	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}

	prodMap := make(map[int] Product)

	for i := 0; i< len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}

	return prodMap, nil
}

func getProduct(productID int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[productID]; ok {
		return &product
	}
	return nil
}

func removeProduct(productID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.m))

	for _, value := range productMap.m {
		products = append(products, value)
	}

	productMap.RUnlock()
	return products
}

func getSortedProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key:= range productMap.m {
		productIds = append(productIds, key)
	}

	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductId() int {
	productIDs := getSortedProductIds()
	return productIDs[len(productIDs) -1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	//if product is set, update it otherwise add it
	addOrUpdateID := -1

	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		if oldProduct == nil {
			return 0, fmt.Errorf("product id: [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductId()
		product.ProductID = addOrUpdateID
	}

	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}