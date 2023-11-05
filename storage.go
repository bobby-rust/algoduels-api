package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	// Account CRUD
	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	// front end will not delete any of these, so i have chosen to omit delete routes
	// Problem CRU - no need for delete 
	CreateProblem(*Problem) error
	GetProblemByID(int) (*Problem, error)
	GetProblems() ([]*Problem, error)
	UpdateProblem(*Problem) error 

	// TestCase CRU - no need for delete
	CreateTestCase(*TestCase) error
	GetTestCaseByID(int) (*TestCase, error)
	GetTestCases() ([]*TestCase, error)
	UpdateTestCase(*TestCase) error 
	// Submission CRU - no need for delete 
	CreateSubmission(*Submission) error
	GetSubmissionByID(int) (*Submission, error)
	GetSubmissions() ([]*Submission, error)
	UpdateSubmissions(*Submission) error

}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	connStr := fmt.Sprintf("host=172.24.97.50 user=%s dbname=%s password=%s sslmode=disable", dbUser, dbName, dbPass)
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

func (s* PostgresStore) createProblemTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Problem (
			problem_id SERIAL PRIMARY KEY,
			prompt VARCHAR(255),
			starter_code TEXT,
			difficulty SMALLINT
		)
	`

	_, err := s.db.Exec(query);
	return err;
}

func (s* PostgresStore) createTestCaseTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS TestCase (
			test_id SERIAL PRIMARY KEY,
			problem_id INT REFERENCES Problem(problem_id),
			input TEXT,
			output TEXT
		)
	`

	_, err := s.db.Exec(query);
	return err;
}

func (s* PostgresStore) createSubmissionTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Submission (
			submission_id SERIAL PRIMARY KEY,
			user_id INT REFERENCES Account(user_id),
			problem_id INT REFERENCES Problem(problem_id),
			submitted_at TIMESTAMP DEFAULT NOW(),
			code TEXT,
			language INT,
			is_accepted BOOLEAN,
			exec_time_ms INT,
			mem_usage_kb INT
		)
	`

	_, err := s.db.Exec(query);
	return err;	
}

// -- Account Create -- 
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

// -- Account Read -- 
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

// -- Account Update -- 
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

// -- Account Delete --
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

// --  Problem Create -- 
func (s *PostgresStore) CreateProblem(acc *Account) error {
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

// -- Problem Read -- 
func (s *PostgresStore) GetProblemByID(id int) (*Account, error) {
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

func (s *PostgresStore) GetProblems() ([]*Account, error) {
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

// -- Problem Update -- 
func (s *PostgresStore) UpdateProblem(*Account) error {
	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.UserID, &account.FirstName, &account.LastName, &account.Username, &account.Email, &account.CreatedAt)

	return account, err
}