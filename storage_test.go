package main

import (
	"log"
	"os"
	"testing"
	"time"
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testStore *PostgresStore

func TestMain(m *testing.M) {
	// Set up
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	testStore, err = NewPostgresStore()
	if err != nil {
		log.Fatalf("Error creating test store: %v", err)
	}

	if err := testStore.Init(); err != nil {
		log.Fatalf("Error initializing test database: %v", err)
	}

	// Run tests
	code := m.Run()

	// Tear down
	if err := testStore.DropTable(); err != nil {
		log.Printf("Error dropping test table: %v", err)
	}

	os.Exit(code)
}

func TestCreateAndGetAccount(t *testing.T) {
	acc := &Account{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@example.com",
		EncryptedPassword: "encrypted_password",
		Phone:             1234567890,
		Balance:           1000,
		CreatedAt:         time.Now().UTC(),
	}

	err := testStore.CreateAccount(acc)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, acc.ID)

	fetchedAcc, err := testStore.GetAccountByID(acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.Email, fetchedAcc.Email)

	fetchedByEmail, err := testStore.GetAccountByEmail(acc.Email)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, fetchedByEmail.ID)
}

func TestGetAccounts(t *testing.T) {
	// Create a few accounts
	for i := 0; i < 3; i++ {
        acc := &Account{
            FirstName:         "Test",
            LastName:          "User",
            Email:             fmt.Sprintf("test%d@example.com", i),
            EncryptedPassword: "password",
            Phone:             1234567890,
            Balance:           1000,
            CreatedAt:         time.Now().UTC(),
        }
        err := testStore.CreateAccount(acc)
        assert.NoError(t, err)
    }

	accounts, err := testStore.GetAccounts()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(accounts), 3)
}

func TestDeleteAccount(t *testing.T) {
	acc := &Account{
		FirstName:         "Jane",
		LastName:          "Doe",
		Email:             "jane@example.com",
		EncryptedPassword: "password",
		Phone:             9876543210,
		Balance:           2000,
		CreatedAt:         time.Now().UTC(),
	}

	err := testStore.CreateAccount(acc)
	assert.NoError(t, err)

	err = testStore.DeleteAccount(acc.ID)
	assert.NoError(t, err)

	_, err = testStore.GetAccountByID(acc.ID)
	assert.Error(t, err)
}

func TestUpdateAccount(t *testing.T) {
	acc := &Account{
		FirstName:         "Alice",
		LastName:          "Smith",
		Email:             "alice@example.com",
		EncryptedPassword: "password",
		Phone:             1122334455,
		Balance:           3000,
		CreatedAt:         time.Now().UTC(),
	}

	err := testStore.CreateAccount(acc)
	assert.NoError(t, err)

	acc.Balance = 3500
    err = testStore.UpdateAccount(acc)
    assert.NoError(t, err)

    updatedAcc, err := testStore.GetAccountByID(acc.ID)
    assert.NoError(t, err)
    assert.Equal(t, int64(3500), updatedAcc.Balance)
}