package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host         = "localhost"
	port         = 5433
	databaseName = "mydatabase"
	username     = "myuser"
	password     = "mypassword"
)

var db *sql.DB

type Product struct {
	ID    int
	Name  string
	Price int
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, databaseName)

	sdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	db = sdb
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// Connection database successful
	fmt.Println("Connection database successful")

	products, err := getProducts()
	fmt.Println("Get Successful !", products)
	if err != nil {
		log.Fatal(err)
	}

	// err = createProduct(&Product{Name: "Product 6", Price: 2000})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// product, err := getProduct(2)
	// fmt.Println("Get Successful !", product)

	// err = updateProduct(4, &Product{Name: "XYZ", Price: 300})
	// product, err := updateProduct(4, &Product{Name: "UUU", Price: 300})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Update Successful !", product)

	// err = deleteProduct(4)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Delete product Successful !")
}

func createProduct(product *Product) error {

	_, err := db.Exec(
		"INSERT INTO public.products (name, price) VALUES ($1, $2);", product.Name, product.Price)

	return err
}

func getProduct(id int) (Product, error) {
	var p Product
	row := db.QueryRow(
		"SELECT id,name,price FROM products WHERE id = $1", id)
	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		fmt.Println("Error Scan", err)
		return Product{}, err
	}

	return p, nil
}

func getProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")

	if err != nil {
		return nil, err
	}

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			fmt.Println("Error Scan", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func updateProduct(id int, product *Product) (Product, error) {
	// _, err := db.Exec(
	// 	"UPDATE public.products	SET name=$1, price =$2 WHERE id = $3;",
	// 	product.Name,
	// 	product.Price,
	// 	id)

	// return err

	var p Product
	row := db.QueryRow("UPDATE public.products	SET name=$1, price =$2 WHERE id = $3 RETURNING id,name,price;", product.Name, product.Price, id)

	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		fmt.Println("Error Scan", err)
		return Product{}, err
	}
	return p, err
}

func deleteProduct(id int) error {
	_, err := db.Exec(
		"DELETE FROM products WHERE id = $1;", id,
	)

	return err
}
