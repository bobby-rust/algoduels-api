package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "host=172.24.97.50 user=postgres dbname=postgres password=test sslmode=disable"
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

func (s* PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s* PostgresStore) createAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Account (
			user_id SERIAL PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			username VARCHAR(50),
			email VARCHAR(50),
			encrypted_password VARCHAR(50),
			created_at TIMESTAMP
		)
	`

	_, err := s.db.Exec(query);
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
			INSERT INTO Account (
				first_name,
				last_name,
				username,
				email,
				encryped_password,
				created_at
			) 
			VALUES ($1, $2, $3, $4, $5, $6)
		`
	
	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Username, acc.Email, acc.EncryptedPassword, acc.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `
		DELETE FROM Account WHERE ID=$1
	`

	_, err := s.db.Query(query, id)
	if err != nil {
		return err
	}

	return nil	
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM Account WHERE ID=$1`

	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	
	return nil, fmt.Errorf("Account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM Account`

	rows, err := s.db.Query(query)

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

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.UserID, &account.FirstName, &account.LastName, &account.Username, &account.Email, &account.CreatedAt)

	return account, err
}