package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
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

	/* --- Front end will not delete any of these below, so i have chosen to omit delete routes --- */

	// Problem CRU - no need for delete
	CreateProblem(*Problem) (int, error)
	GetProblemByID(int) (*Problem, error)
	GetProblems() ([]*Problem, error)
	UpdateProblem(*Problem) error

	// TestCase CRU - no need for delete
	CreateTestCase(*TestCase) (int, error)
	GetTestCaseByProblemID(int) ([]*TestCase, error)
	GetTestCases() ([]*TestCase, error)
	UpdateTestCase(*TestCase) error

	// Submission CRU - no need for delete
	CreateSubmission(*Submission) error
	GetSubmissionByID(int) (*Submission, error)
	GetSubmissions() ([]*Submission, error)
	UpdateSubmission(*Submission) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	/* Load environment variables */
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbIP := os.Getenv("JUDGE0_IP")

	/* Create connection string */
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", dbIP, dbUser, dbName, dbPass)

	/* Open database connection */
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Database connection opened")
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	tableCreationFuncs := []func() error{
		s.createAccountTable,
		s.createProblemTable,
		s.createTestCaseTable,
		s.createSubmissionTable,
	}

	for _, f := range tableCreationFuncs {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) createAccountTable() error {
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

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createProblemTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Problem (
			problem_id SERIAL PRIMARY KEY,
			prompt VARCHAR(255),
			starter_code TEXT,
			difficulty SMALLINT
		)
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createTestCaseTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS TestCase (
			test_case_id SERIAL PRIMARY KEY,
			problem_id INT REFERENCES Problem(problem_id),
			input TEXT,
			output TEXT,
            is_sanity_check BOOLEAN
		)
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createSubmissionTable() error {
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

	_, err := s.db.Exec(query)
	return err
}

// -- Account Create --
func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
			INSERT INTO Account (
				username,
				email,
				encrypted_password,
				created_at
			) 
			VALUES ($1, $2, $3, $4)
		`

	_, err := s.db.Query(query, acc.Username, acc.Email, acc.EncryptedPassword, acc.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// -- Account Read --
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM Account WHERE user_id=$1`

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
		DELETE FROM Account WHERE user_id=$1
	`

	_, err := s.db.Query(query, id)
	if err != nil {
		return err
	}

	return nil
}

// --  Problem Create --
func (s *PostgresStore) CreateProblem(prob *Problem) (int, error) {
	query := `
			INSERT INTO Problem (
				prompt,
				starter_code,
				difficulty
			) 
			VALUES ($1, $2, $3) RETURNING problem_id;
		`
	var problemID int
	err := s.db.QueryRow(query, prob.Prompt, prob.StarterCode, prob.Difficulty).Scan(&problemID)
	fmt.Printf("ProblemID: %d", problemID)
	if err != nil {
		return -1, err // -1 signifies an error occurred
	}

	return problemID, nil
}

// -- Problem Read --
func (s *PostgresStore) GetProblemByID(id int) (*Problem, error) {
	query := `SELECT * FROM Problem WHERE problem_id=$1`

	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoProblem(rows)
	}

	return nil, fmt.Errorf("Account %d not found", id)
}

func (s *PostgresStore) GetProblems() ([]*Problem, error) {
	query := `SELECT * FROM Problem`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	problems := []*Problem{}

	for rows.Next() {
		problem, err := scanIntoProblem(rows)
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

// -- Problem Update --
func (s *PostgresStore) UpdateProblem(*Problem) error {
	return nil
}

// --  TestCase Create --
func (s *PostgresStore) CreateTestCase(testcase *TestCase) (int, error) {
	query := `
			INSERT INTO TestCase (
				problem_id, 
				input,
				output,
                is_sanity_check
			) 
			VALUES ($1, $2, $3, $4) RETURNING test_case_id;
		`

	var testCaseID int
	err := s.db.QueryRow(query, testcase.ProblemID, testcase.Input, testcase.Output, testcase.IsSanityCheck).Scan(&testCaseID)
	if err != nil {
		return -1, err
	}

	return testCaseID, nil
}

// -- TestCase Read -- ID here is a PROBLEM id
func (s *PostgresStore) GetTestCaseByProblemID(id int) ([]*TestCase, error) {
	query := `SELECT * FROM TestCase WHERE problem_id=$1`

	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	testCases := []*TestCase{}

	for rows.Next() {
		testCase, err := scanIntoTestCase(rows)
		if err != nil {
			return nil, err
		}
		testCases = append(testCases, testCase)
	}

	return testCases, nil
}

func (s *PostgresStore) GetTestCases() ([]*TestCase, error) {
	query := `SELECT * FROM TestCase`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	testCases := []*TestCase{}

	for rows.Next() {
		testCase, err := scanIntoTestCase(rows)
		if err != nil {
			return nil, err
		}
		testCases = append(testCases, testCase)
	}

	return testCases, nil
}

// -- TestCase Update --
func (s *PostgresStore) UpdateTestCase(*TestCase) error {
	return nil
}

// --  Problem Create --
func (s *PostgresStore) CreateSubmission(sub *Submission) error {
	query := `
			INSERT INTO Submission (
				user_id,
				problem_id,
				submitted_at,
				code,
				language,
				is_accepted,
				exec_time_ms,
				mem_usage_kb
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

	_, err := s.db.Query(query, sub.UserID, sub.ProblemID, sub.SubmittedAt, sub.Code, sub.Language, sub.IsAccepted, sub.ExecTimeMS, sub.MemUsageKB)
	if err != nil {
		return err
	}

	return nil
}

// -- Problem Read --
func (s *PostgresStore) GetSubmissionByID(id int) (*Submission, error) {
	query := `SELECT * FROM Submission WHERE submission_id=$1`

	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoSubmission(rows)
	}

	return nil, fmt.Errorf("Submission %d not found", id)
}

func (s *PostgresStore) GetSubmissions() ([]*Submission, error) {
	query := `SELECT * FROM Submission`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	subs := []*Submission{}

	for rows.Next() {
		sub, err := scanIntoSubmission(rows)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

// -- Problem Update --
func (s *PostgresStore) UpdateSubmission(*Submission) error {
	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.UserID, &account.Username, &account.Email, &account.CreatedAt)

	return account, err
}

func scanIntoSubmission(rows *sql.Rows) (*Submission, error) {
	sub := new(Submission)
	err := rows.Scan(&sub.UserID, &sub.ProblemID, &sub.SubmittedAt, &sub.Code, &sub.Language, &sub.IsAccepted, &sub.ExecTimeMS, &sub.MemUsageKB)

	return sub, err
}

func scanIntoTestCase(rows *sql.Rows) (*TestCase, error) {
	tc := new(TestCase)
	err := rows.Scan(&tc.TestCaseID, &tc.ProblemID, &tc.Input, &tc.Output, &tc.IsSanityCheck)

	return tc, err
}

func scanIntoProblem(rows *sql.Rows) (*Problem, error) {
	p := new(Problem)
	err := rows.Scan(&p.Prompt, &p.StarterCode, &p.Difficulty)

	return p, err
}
