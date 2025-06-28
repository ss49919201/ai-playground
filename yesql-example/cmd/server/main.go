package main

import (
	"fmt"
	"log"
	"net/http"
	"yesql-account-system/internal/account"
	"yesql-account-system/internal/auth"
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
	authService := auth.NewService(database.GetConn(), loader)
	handler := NewHandler(accountService, authService)

	authMiddleware := auth.NewMiddleware(authService)

	r := mux.NewRouter()
	
	// Auth routes (no authentication required)
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	
	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	
	protected.HandleFunc("/logout", handler.Logout).Methods("POST")
	protected.HandleFunc("/accounts", handler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts", handler.GetUserAccounts).Methods("GET")
	protected.HandleFunc("/accounts/{id}", handler.GetAccount).Methods("GET")
	protected.HandleFunc("/accounts/{id}/transactions", handler.GetAccountTransactions).Methods("GET")
	protected.HandleFunc("/deposit", handler.Deposit).Methods("POST")
	protected.HandleFunc("/withdraw", handler.Withdraw).Methods("POST")
	protected.HandleFunc("/transfer", handler.Transfer).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}