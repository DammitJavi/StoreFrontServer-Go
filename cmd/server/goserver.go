package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "javierrojas"
	dbname = "storefront"
)

type Inventory struct{
	id 				int64
	product_name 	string
	category 		string
	price 			float32
	quantity		int64
	sku				string
	barcode			int64
	supplier		string
	last_restock_date  	string
	low_restock_threshold int64
	weight				float32
	dimensions		string
	status			string
}

func handler(w http.ResponseWriter, r *http.Request){

	var invs []Inventory

	// fmt.Fprintf(w, "Hello %s", r.URL.Path[1:])


	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " + "dbname=%s sslmode=disable", host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()


	if err != nil {
		panic(err)
	}

	fmt.Println("Postgres DB connected.")


	rows, err := db.Query("SELECT * FROM Inventory")
	
	if err != nil {
		fmt.Errorf("Error: %v", err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next(){
		var inv Inventory
		if err := rows.Scan(&inv.id, &inv.product_name, &inv.category, &inv.price, &inv.quantity, &inv.sku, &inv.barcode, &inv.supplier, &inv.last_restock_date, &inv.low_restock_threshold, &inv.weight, &inv.dimensions, &inv.status ); err != nil {
			fmt.Errorf("Error getting from db: %v", err)
			panic(err)
		}
		invs = append(invs, inv)
	}

	if err := rows.Err(); err != nil {
		fmt.Errorf("Error at end: %v", err)
		panic(err)
	}


	fmt.Printf("This is inventory from db: %v\n", invs)
}


func main(){

	http.HandleFunc("/", handler)
	fmt.Println("Listening on Port 3000")
	err := http.ListenAndServe(":3000", nil)

	if err != nil{
		panic(err)
	}

}	
