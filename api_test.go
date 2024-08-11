package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateAccount(account *Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockStorage) DeleteAccount(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorage) UpdateAccount(account *Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockStorage) GetAccounts() ([]*Account, error) {
	args := m.Called()
	return args.Get(0).([]*Account), args.Error(1)
}

func (m *MockStorage) GetAccountByID(id int) (*Account, error) {
	args := m.Called(id)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockStorage) GetAccountByEmail(email string) (*Account, error) {
	args := m.Called(email)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockStorage) DropTable() error {
	args := m.Called()
	return args.Error(0)
}

func TestHandleAccount(t *testing.T) {
	mockStorage := new(MockStorage)
	server := NewAPIServer(":8080", mockStorage)

	t.Run("Create Account", func(t *testing.T) {
		newAccount := &Account{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Phone:     1234567890,
		}

		mockStorage.On("CreateAccount", mock.AnythingOfType("*main.Account")).Return(nil)

		body, _ := json.Marshal(newAccount)
		req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		server.handleCreateAccount(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Get Account", func(t *testing.T) {
		account := &Account{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Phone:     1234567890,
		}

		mockStorage.On("GetAccountByID", 1).Return(account, nil)

		req, _ := http.NewRequest("GET", "/account/1", nil)
		rr := httptest.NewRecorder()

		server.handleGetAccountByID(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockStorage.AssertExpectations(t)
	})
}

func TestHandleGetAccounts(t *testing.T) {
	mockStorage := new(MockStorage)
	server := NewAPIServer(":8080", mockStorage)

	accounts := []*Account{
		{ID: 1, FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		{ID: 2, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
	}

	mockStorage.On("GetAccounts").Return(accounts, nil)

	req, _ := http.NewRequest("GET", "/accounts", nil)
	rr := httptest.NewRecorder()

	server.handleGetAllAccounts(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockStorage.AssertExpectations(t)

	var responseAccounts []*Account
	json.Unmarshal(rr.Body.Bytes(), &responseAccounts)
	assert.Equal(t, accounts, responseAccounts)
}
