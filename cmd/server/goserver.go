package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

const (
	host = "localhost"
	port = 5432
	user = "javierrojas"
	dbname = "storefront"

	emailRegex = `^[a-zA-Z0-9.+-_]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordRegex = `^[a-zA-Z0-9!#@%\$\&]+$`
)

// type Inventory struct{
// 	id 				int64
// 	product_name 	string
// 	category 		string
// 	price 			float32
// 	quantity		int64
// 	sku				string
// 	barcode			int64
// 	supplier		string
// 	last_restock_date  	string
// 	low_restock_threshold int64
// 	weight				float32
// 	dimensions		string
// 	status			string
// }

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

type MainPageInv struct {
	ID 				int64	`json:"id"`
	Product_name 	string	`json:"product_name"`
	Category 		string	`json:"category"`
	Price 			float32	`json:"price"`
	Sku				string	`json:"sku"`
	Dimensions		string	`json:"dimensions"`
	Status			string	`json:"status"`
}

type productById struct {
	ID 				int64	`json:"id"`
	Product_name 	string	`json:"product_name"`
	Category 		string	`json:"category"`
	Price 			float32	`json:"price"`
	Sku				string	`json:"sku"`
	Supplier 		string 	`json:"supplier"`
	Dimensions		string	`json:"dimensions"`
	Status			string	`json:"status"`
}

type users struct {
	Username	string `json:"username"`
	Email 		string `json:"email"`
	Password  	string `json:"password"`
}

func mainPageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

		if r.Method == http.MethodPost {
			http.Error(w, "Method not allowed...", http.StatusMethodNotAllowed)
			return
		}

		var mainInv []MainPageInv

		rows, err := db.Query("SELECT id, product_name, category, price, sku, dimensions, status FROM Inventory")
	
		
		if err != nil {
			fmt.Errorf("Error retrieving data from database; %v", err)
			return 
		}

		defer rows.Close()
		
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

		// fmt.Println(string(jData))
		if err != nil {
			fmt.Errorf("Error with jData: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jData)

	}
}

func productByIdHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		//Post
		if r.Method == http.MethodPost {
			var keys []int
			var products []productById

			err := json.NewDecoder(r.Body).Decode(&keys)

			if err != nil {
				fmt.Errorf("Error in post for cart items: %v", err)
				return
			}

			rows, err := db.Query("SELECT id, product_name, category, price, sku, supplier, dimensions, status from Inventory where id = ANY($1)", pq.Array(keys) )
			
			if err != nil {
				fmt.Errorf("Error with query for keys: %v", err)
				return
			}

			defer rows.Close()

			for rows.Next(){
				var product productById

				if err := rows.Scan(&product.ID, &product.Product_name, &product.Category, &product.Price, &product.Sku, &product.Supplier, &product.Dimensions, & product.Status); err != nil {
					fmt.Errorf("Error with single item in keys: %v", err)
				}
				products = append(products, product)

			}
			jData, err := json.Marshal(products)

			if err != nil {
				fmt.Errorf("Json Marshal error: %v", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jData)

		}else{
			//Get
			var product2 productById
			path := strings.TrimPrefix(r.URL.Path, "/api/product/")
			id := strings.SplitN(path, "/", 2)[0]
	
			rows := db.QueryRow("SELECT id, product_name, category, price, sku, supplier, dimensions, status from Inventory where id = $1", id)
			err := rows.Scan( &product2.ID, &product2.Product_name, &product2.Category, &product2.Price, &product2.Sku, &product2.Supplier, &product2.Dimensions, & product2.Status )
			// fmt.Println(product)

			if err != nil {
				fmt.Errorf("Error while getting product by ID: %v", err)
			}
			
			jData, err := json.Marshal(product2)
	
			if err != nil {
				fmt.Errorf("Error while converting product by id into JSON: %v", err)
	
			}
	
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jData)
		}
	}
}

func userHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			log.Println("Error with user, No GET allowed.")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user users

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, "User info error", http.StatusBadRequest)
			return
		}

		emailRe := regexp.MustCompile(emailRegex)
		
		if !emailRe.MatchString(user.Email) {
			log.Println("Email did not meet regex.")
			http.Error(w, "Email did not meet regex.", http.StatusBadRequest)
			return
		}
		
		passwordRe := regexp.MustCompile(passwordRegex)

		if !passwordRe.MatchString(user.Password){
			log.Println("Password did not meet regex.")
			http.Error(w, "Password did not meet regex.", http.StatusBadRequest)
			return
		}

		securePassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Println("Bcrypting password did not work")
			http.Error(w, "Bcrypt Error.", http.StatusExpectationFailed)
			return
		}

		_, err2 := db.Exec("INSERT INTO usersdb(username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, securePassword)

		if err2 != nil {
			log.Println("Error with insert user into db", err2)
			http.Error(w, "User could not be added to db", http.StatusBadRequest)
			return
		}

		fmt.Printf("User %s has been added.", user.Username)
	}
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			log.Println("Method Get is not allowed")
			http.Error(w, "Method Get is not allowed.", http.StatusBadRequest)
			return
		}

		var user users
		var userQ users

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			log.Printf("Error with request body: %v \n", err)
			http.Error(w, "Error with request body", http.StatusBadRequest)
			return
		}

		rows := db.QueryRow("SELECT username, password FROM usersdb WHERE username = $1", user.Username)
		
		err2 := rows.Scan(&userQ.Username, &userQ.Password)

		if err2 != nil {
			log.Println("Problem with db query.")
			http.Error(w, "Problem with db query.", http.StatusBadRequest)
			return
		}

		err3 := bcrypt.CompareHashAndPassword([]byte(userQ.Password), []byte(user.Password))

		if err3 != nil {
			log.Println("Wrong username or password.")
			http.Error(w, "Wrong username or password.", http.StatusBadRequest)
			return
		}

		fmt.Printf("Welcome %s!", userQ.Username)

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

	// mux := http.NewServeMux()

	// Bottom code is to test database
	// http.HandleFunc("/", handler)  

	http.HandleFunc("/api/", mainPageHandler(db))
	http.HandleFunc("/api/product/", productByIdHandler(db))
	http.HandleFunc("/api/users/", userHandler(db))
	http.HandleFunc("/api/login/", loginHandler(db))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(http.DefaultServeMux)
	fmt.Println("Listening on Port 3000")
	http.ListenAndServe(":3000", handler)
	
}	
