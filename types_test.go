package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateNewAccount(t *testing.T) {
	testCases := []struct {
		name      string
		firstName string
		lastName  string
		email     string
		password  string
	}{
		{"Valid account", "John", "Doe", "john@example.com", "password123"},
		{"Empty password", "Jane", "Smith", "jane@example.com", ""},
		{"Empty email", "Bob", "Johnson", "", "securepass"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			acc, err := GenerateNewAccount(tc.firstName, tc.lastName, tc.email, tc.password)

			if tc.password == "" || tc.email == "" {
				assert.Error(t, err)
				assert.Nil(t, acc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, acc)
				assert.Equal(t, tc.firstName, acc.FirstName)
				assert.Equal(t, tc.lastName, acc.LastName)
				assert.Equal(t, tc.email, acc.Email)
				assert.NotEmpty(t, acc.EncryptedPassword)
				assert.NotZero(t, acc.Phone)
				assert.WithinDuration(t, time.Now().UTC(), acc.CreatedAt, 2*time.Second)

				// Verify password encryption
				err = bcrypt.CompareHashAndPassword([]byte(acc.EncryptedPassword), []byte(tc.password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoginRequest(t *testing.T) {
	lr := LoginRequest{
		Email:             "test@example.com",
		EncryptedPassword: "hashedpassword",
	}

	assert.Equal(t, "test@example.com", lr.Email)
	assert.Equal(t, "hashedpassword", lr.EncryptedPassword)
}

func TestTransferRequest(t *testing.T) {
	tr := TransferRequest{
		ToAccount: 123,
		Amount:    1000,
	}

	assert.Equal(t, 123, tr.ToAccount)
	assert.Equal(t, 1000, tr.Amount)
}

func TestAccount(t *testing.T) {
	now := time.Now().UTC()
	acc := Account{
		ID:                1,
		FirstName:         "Alice",
		LastName:          "Wonder",
		Email:             "alice@wonderland.com",
		Phone:             1234567890,
		EncryptedPassword: "hashedpassword",
		Balance:           5000,
		CreatedAt:         now,
	}

	assert.Equal(t, 1, acc.ID)
	assert.Equal(t, "Alice", acc.FirstName)
	assert.Equal(t, "Wonder", acc.LastName)
	assert.Equal(t, "alice@wonderland.com", acc.Email)
	assert.Equal(t, int64(1234567890), acc.Phone)
	assert.Equal(t, "hashedpassword", acc.EncryptedPassword)
	assert.Equal(t, int64(5000), acc.Balance)
	assert.Equal(t, now, acc.CreatedAt)
}

func TestNewAccount(t *testing.T) {
	na := NewAccount{
		FirstName:         "Charlie",
		LastName:          "Brown",
		Email:             "charlie@peanuts.com",
		EncryptedPassword: "snoopy",
	}

	assert.Equal(t, "Charlie", na.FirstName)
	assert.Equal(t, "Brown", na.LastName)
	assert.Equal(t, "charlie@peanuts.com", na.Email)
	assert.Equal(t, "snoopy", na.EncryptedPassword)
}