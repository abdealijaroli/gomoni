package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store      Storage 
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /account", makeHTTPHandleFunc(s.handleGetAllAccounts))
	router.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))
	router.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	log.Println("API server running on port:", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, router); err != nil {
		log.Fatal("Error booting up the server:", err)
	}
}

func (s *APIServer) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return fmt.Errorf("invalid id")
	}

	res, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, res)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	newAccount := &NewAccount{}
	if err := json.NewDecoder(r.Body).Decode(newAccount); err != nil {
		return err
	}
	account := GenerateNewAccount(newAccount.FirstName, newAccount.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleFundTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
