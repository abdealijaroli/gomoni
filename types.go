package main

import (
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email             string `json:"email"`
	EncryptedPassword string `json:"password"`
}

type TransferRequest struct {
	FromAccount int   `json:"fromAccount"`
	ToAccount   int   `json:"toAccount"`
	Amount      int64 `json:"amount"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Email             string    `json:"email"`
	Phone             int64     `json:"phone"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

type NewAccount struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
}

func GenerateNewAccount(firstname, lastname, email, password string) (*Account, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	enpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstname,
		LastName:          lastname,
		Email:             email,
		EncryptedPassword: string(enpw),
		Phone:             int64(rand.Intn(1e5)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
