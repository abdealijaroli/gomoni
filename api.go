package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	router.HandleFunc("POST /login", makeHTTPHandleFunc(s.handleLogin, false))
	router.HandleFunc("GET /account", authWithJWT(makeHTTPHandleFunc(s.handleGetAllAccounts, true), s.store))
	router.HandleFunc("POST /account", authWithJWT(makeHTTPHandleFunc(s.handleCreateAccount, true), s.store))
	router.HandleFunc("GET /account/{id}", authWithJWT(makeHTTPHandleFunc(s.handleGetAccountByID, true), s.store))
	router.HandleFunc("DELETE /account/{id}", authWithJWT(makeHTTPHandleFunc(s.handleDeleteAccount, true), s.store))
	router.HandleFunc("POST /transfer", authWithJWT(makeHTTPHandleFunc(s.handleTransfer, true), s.store))

	log.Println("API server running on port:", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, router); err != nil {
		log.Fatal("Error booting up the server:", err)
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByEmail(loginReq.Email)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, APIError{Error: "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.EncryptedPassword), []byte(loginReq.EncryptedPassword)); err != nil {
		return WriteJSON(w, http.StatusUnauthorized, APIError{Error: "Invalid credentials"})
	}

	tokenString, err := createJWT(acc)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   false, // true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	return WriteJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "Login successful"})
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
	account, err := GenerateNewAccount(newAccount.FirstName, newAccount.LastName, newAccount.Email, newAccount.EncryptedPassword)
	if err != nil {
		return err
	}
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

	fromAccount, err := s.store.GetAccountByID(transferReq.FromAccount)
	if err != nil {
		return err
	}

	if fromAccount.Balance < transferReq.Amount {
		return fmt.Errorf("insufficient funds")
	}

	toAccount, err := s.store.GetAccountByID(transferReq.ToAccount)
	if err != nil {
		return err
	}

	fromAccount.Balance -= transferReq.Amount
	toAccount.Balance += transferReq.Amount

	if err := s.store.UpdateAccount(fromAccount); err != nil {
		return err
	}

	if err := s.store.UpdateAccount(toAccount); err != nil {
		fromAccount.Balance += transferReq.Amount
		s.store.UpdateAccount(fromAccount)
		return err
	}

	return WriteJSON(w, http.StatusOK, transferReq)
}

func unauthorized(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, APIError{Error: "Unauthorized"})
}

func createJWT(account *Account) (string, error) {
	claims := jwt.MapClaims{
		"id":    account.ID,
		"email": account.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func authWithJWT(f http.HandlerFunc, store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if cookie, err := r.Cookie("token"); err == nil {
			tokenString = cookie.Value
		}
		if tokenString == "" {
			unauthorized(w) // Token not found
			return
		}

		token, err := validateJWT(tokenString)
		if err != nil || !token.Valid {
			unauthorized(w) // Invalid token
			return
		}

		// Rest of the function remains the same
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthorized(w) // Invalid claims
			return
		}

		accountID, ok := claims["id"].(float64)
		if !ok {
			unauthorized(w) // Invalid account ID
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			unauthorized(w) // Invalid email
			return
		}

		account, err := store.GetAccountByID(int(accountID))
		if err != nil {
			unauthorized(w) // User not found
			return
		}

		if account.Email != email {
			unauthorized(w) // Email mismatch
			return
		}

		ctx := NewAuthContext(r.Context(), int(accountID), email)
		r = r.WithContext(ctx)

		f.ServeHTTP(w, r)
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

func makeHTTPHandleFunc(f APIFunc, requireAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if requireAuth {
			authCtx, ok := GetAuthContext(r.Context())
			if !ok {
				unauthorized(w) // Unauthorized
				return
			}
			log.Printf("Authenticated request: AccountID=%d, Email=%s\n", authCtx.AccountID, authCtx.Email)
		}

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
