package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      storage
}

func NewAPIService(listenAddr string, store storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))

	// listen and serve
	log.Println("listening on port", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, router); err != nil {
		log.Println("Error listening to the server!")
		os.Exit(1)
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, _ *http.Request) error {
	accounts, getAccountsErr := s.store.GetAccounts()
	if getAccountsErr != nil {
		log.Println("error while getting accounts.")
		return getAccountsErr
	}

	log.Println("Received accounts successfully.", accounts)
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccReq := CreateAccountRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(&createAccReq); decodeErr != nil {
		return decodeErr
	}

	newAccount := newAccount(createAccReq.FirstName, createAccReq.LastName)

	insertedId, createAccErr := s.store.CreateAccount(newAccount)
	if createAccErr != nil {
		log.Println("Error while creating a new account!")
		return createAccErr
	}
	log.Println("Account created successfully")
	newAccount.ID = int(insertedId)
	return WriteJSON(w, http.StatusOK, newAccount)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

// wrapper function: Error Handler
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
