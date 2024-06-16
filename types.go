package main

import "math/rand"

type AccountDetails struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     int64  `json:"phone"`
	Balance   int64  `json:"balance"`
}

func CreateNewAccount(firstname, lastname string) *AccountDetails {
	return &AccountDetails{
		ID:        rand.Intn(1e5),
		FirstName: firstname,
		LastName:  lastname,
		Phone:     int64(rand.Intn(10*1e9)),
		Balance:   int64(rand.Intn(1e5)),
	}
}
