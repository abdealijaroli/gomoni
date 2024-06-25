package main

import (
	"math/rand"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     int64     `json:"phone"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewAccount struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func GenerateNewAccount(firstname, lastname string) *Account {
	return &Account{
		FirstName: firstname,
		LastName:  lastname,
		Phone:     int64(rand.Intn(1e5)),
		CreatedAt: time.Now().UTC(),
	}
}
