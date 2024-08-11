package main

import (
	"database/sql"
	"fmt"
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
	GetAccountByEmail(string) (*Account, error)
	DropTable() error
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) DropTable() error {
	_, err := s.db.Exec("DROP TABLE account")
	return err
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
		email varchar(50),
		encrypted_password varchar(100),
		phone bigint,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	rows, err := s.db.Query(`select * from account where email=$1`, email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %s not found", email)
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	q := `insert into 
		account(first_name, last_name, email, encrypted_password, phone, balance, created_at)
		values($1, $2, $3, $4, $5, $6, $7)
		returning id
	`
	err := s.db.QueryRow(q, acc.FirstName, acc.LastName, acc.Email, acc.EncryptedPassword, acc.Phone, acc.Balance, acc.CreatedAt).Scan(&acc.ID)

	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	q := `delete from account where id=$1`

	_, err := s.db.Query(q, id)
	return err
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	q := `UPDATE account SET first_name=$1, last_name=$2, email=$3, encrypted_password=$4, phone=$5, balance=$6 WHERE id=$7`

	_, err := s.db.Exec(q, account.FirstName, account.LastName, account.Email, account.EncryptedPassword, account.Phone, account.Balance, account.ID)
	if err != nil {
		return fmt.Errorf("error updating account: %v", err)
	}

	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`select * from account`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`select * from account where id=$1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}

	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.EncryptedPassword,
		&account.Phone,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
