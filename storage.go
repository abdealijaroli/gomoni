package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	connStr := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		phone serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	q := `insert into 
		account(first_name, last_name, phone, balance, created_at)
		values($1, $2, $3, $4, $5)
	`
	_, err := s.db.Exec(q, acc.FirstName, acc.LastName, acc.Phone, acc.Balance, acc.CreatedAt)

	return err
}

func (s *PostgresStore) DeleteAccount(int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`select * from account`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account := &Account{}

		err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Phone, &account.Balance, &account.CreatedAt)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(int) (*Account, error) {
	return nil, nil
}
