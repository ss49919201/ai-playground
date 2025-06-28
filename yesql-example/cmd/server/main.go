package main

import (
	"fmt"
	"log"
	"net/http"
	"yesql-account-system/internal/account"
	"yesql-account-system/internal/db"
	"yesql-account-system/internal/yesql"

	"github.com/gorilla/mux"
)

func main() {
	database, err := db.NewDB("accounts.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	err = database.InitSchema("sql/schema.sql")
	if err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	loader := yesql.NewQueryLoader()
	err = loader.LoadQueriesFromDir("sql/queries")
	if err != nil {
		log.Fatalf("Failed to load queries: %v", err)
	}

	accountService := account.NewService(database.GetConn(), loader)
	handler := NewHandler(accountService)

	r := mux.NewRouter()
	
	r.HandleFunc("/accounts", handler.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts", handler.ListAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}", handler.GetAccount).Methods("GET")
	r.HandleFunc("/accounts/{id}/transactions", handler.GetAccountTransactions).Methods("GET")
	r.HandleFunc("/deposit", handler.Deposit).Methods("POST")
	r.HandleFunc("/withdraw", handler.Withdraw).Methods("POST")
	r.HandleFunc("/transfer", handler.Transfer).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}