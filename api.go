package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
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
	router.HandleFunc("GET /account/{id}", authWithJWT(makeHTTPHandleFunc(s.handleGetAccountByID)))
	router.HandleFunc("DELETE /account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount))
	router.HandleFunc("POST /transfer", makeHTTPHandleFunc(s.handleTransfer))

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
	id, err := getID(r)
	if err != nil {
		return err
	}

	acc, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
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
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := &TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func authWithJWT(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("jwt here")
		f(w, r)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return id, fmt.Errorf("invalid ID: %d", id)
	}
	return id, nil
}
