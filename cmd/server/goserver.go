package main

import (
	"database/sql"
	"encoding/json"
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

type MainPageInv struct {
	ID 				int64	`json:"id"`
	Product_name 	string	`json:"product_name"`
	Category 		string	`json:"category"`
	Price 			float32	`json:"price"`
	Sku				string	`json:"sku"`
	Dimensions		string	`json:"dimensions"`
	Status			string	`json:"status"`
}

// func handler(w http.ResponseWriter, r *http.Request){

// 	var invs []Inventory

// 	// fmt.Fprintf(w, "Hello %s", r.URL.Path[1:])

// 	rows, err := db.Query("SELECT * FROM Inventory")
	
// 	if err != nil {
// 		fmt.Errorf("Error: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next(){
// 		var inv Inventory
// 		if err := rows.Scan(&inv.id, &inv.product_name, &inv.category, &inv.price, &inv.quantity, &inv.sku, &inv.barcode, &inv.supplier, &inv.last_restock_date, &inv.low_restock_threshold, &inv.weight, &inv.dimensions, &inv.status ); err != nil {
// 			fmt.Errorf("Error getting from db: %v", err)
// 		}
// 		invs = append(invs, inv)
// 	}

// 	if err := rows.Err(); err != nil {
// 		fmt.Errorf("Error at end: %v", err)
// 	}


// 	fmt.Printf("This is inventory from db: %v\n", invs)
// }

func mainPageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

		if r.Method == http.MethodPost {
			http.Error(w, "Method not allowed...", http.StatusMethodNotAllowed)
			return
		}

		var mainInv []MainPageInv

		rows, err := db.Query("SELECT id, product_name, category, price, sku, dimensions, status FROM Inventory")
		defer rows.Close()

		if err != nil {
			fmt.Errorf("Error retrieving data from database; %v", err)
			return 
		}
		
		for rows.Next(){
			var inv MainPageInv
			if err := rows.Scan(&inv.ID, &inv.Product_name, &inv.Category, &inv.Price, &inv.Sku, &inv.Dimensions, &inv.Status ); err != nil {
				fmt.Errorf("Error getting from db: %v", err)
			}
			mainInv = append(mainInv, inv)
		}

		//Check data
		// fmt.Printf("This is inventory from db: %v\n", mainInv)


		jData, err := json.Marshal(mainInv)


		fmt.Println(string(jData))
		if err != nil {
			fmt.Errorf("Error with jData: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jData)

	}
}


func main(){

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

	// Bottom code is to test database
	// http.HandleFunc("/", handler)  

	http.HandleFunc("/api", mainPageHandler(db))

	fmt.Println("Listening on Port 3000")
	http.ListenAndServe(":3000", nil)
	


}	
