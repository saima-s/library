package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

const (
	API_PATH = "/apis/v1/books"
)

type library struct {
	DbHost     string
	DbPassword string
	DbName     string
}

type Book struct {
	Id         int
	Name, Isbn string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}
	l := library{
		DbHost:     dbHost,
		DbPassword: dbPassword,
		DbName:     dbName,
	}
	r := mux.NewRouter()
	r.HandleFunc("/apis/v1/books", l.getBooks).Methods(http.MethodGet)
	r.HandleFunc("/apis/v1/books", l.addBooks).Methods(http.MethodPost)
	fmt.Println("Starting server!!!")
	http.ListenAndServe(":8080", r)
}

func (l library) addBooks(w http.ResponseWriter, r *http.Request) {
	book := Book{}
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Fatalf("error in unmrshalling the body", err.Error())
	}
	db := l.openConnection()
	defer l.closeConnection(db)
	insertQuery, err := db.Prepare("INSERT INTO books VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("error is prep data : %v", err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("error is begin txn data : %v", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf("error is exec txn data : %v", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("error is commit txn data : %v", err.Error())
	}
	json.NewEncoder(w).Encode("Data added success!!")
	// log.Println("this method is called")
}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	db := l.openConnection()
	defer l.closeConnection(db)
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("error is querying data : %v", err.Error())
	}
	books := []Book{}
	for rows.Next() {
		var id int
		var name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("error in scanning:::%v", err.Error())
		}
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)

	}

	json.NewEncoder(w).Encode(books)
	// log.Println("this method is called")
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.DbPassword, l.DbHost, l.DbName))
	if err != nil {
		log.Fatalf("error in connecting DB with err : %v", err.Error())
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("error in closing connection:%v", err.Error())
	}

}
