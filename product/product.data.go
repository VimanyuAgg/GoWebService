package product

import (
	"context"
	"database/sql"
	"github.com/goWebServices/one-o-one/database"
	"log"
	"time"
)

func getProduct(productID int) (*Product, error) {
	product:= &Product{}
	ctx, cancel := context.WithTimeout(context.Background(), 15* time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products
	WHERE productID = ?`,productID)

	err := row.Scan(&product.ProductID, &product.Manufacturer, &product.Sku, &product.Upc, &product.PricePerUnit,
		&product.QuantityOnHand, &product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Fatal("Error while fetching the details for productID: ",productID)
		return nil, err
	}
	return product, nil
}

func removeProduct(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15* time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products where productId = ?`, productID)
	if err != nil {
		log.Printf("Error occurred while deleting productID: [%d]", productID)
		return err
	}
	return nil
}

func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products
	`)

	if err != nil {
		log.Fatal("Error occured while fetching productList from DB", err)
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID, &product.Manufacturer, &product.Sku, &product.Upc, &product.PricePerUnit,
			&product.QuantityOnHand, &product.ProductName)
		products = append(products, product)
	}

	return products, nil
}


func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET
	manufacturer=?,
	sku=?,
	upc=?,
	pricePerUnit=CAST(? AS DECIMAL(13,2)),
	quantityOnHand=?,
	productName=?
	WHERE productId=?`,
	product.Manufacturer,
	product.Sku,
	product.Upc,
	product.PricePerUnit,
	product.QuantityOnHand,
	product.ProductName,
	product.ProductID)
	if err != nil {
		log.Printf("Error occurred while updating ProductID: %d", product.ProductID)
		log.Println(err)
		return err
	}
	return nil
}

func insertNewProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()
	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products
	(
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName) VALUES (?,?,?,?,?,?)`, product.Manufacturer, product.Sku, product.Upc, product.PricePerUnit, product.QuantityOnHand, product.ProductName)
	if err != nil {
		log.Println("Error occurred while inserting a new product into the DB")
		return 0, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error occurred while getting the lastInsertID of the newly inserted object")
		return 0, nil
	}
	return int(insertID), nil
}
