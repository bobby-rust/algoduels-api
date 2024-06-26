package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type Storage interface {
	// Account CRUD
	CreateAccount(*CreateAccountRequest) (*CreateAccountResponse, error)
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error

	/* --- Front end will not delete any of these below, so i have chosen to omit delete routes --- */

	// Problem CRU - no need for delete
	CreateProblem(*Problem) (int, error)
	GetProblemByID(int) (*Problem, error)
	GetProblemByName(string) (*Problem, error)
	GetProblems() ([]*Problem, error)
	UpdateProblem(*Problem) error

	// TestCase CRU - no need for delete
	CreateTestCase(*TestCase) (int, error)
	GetTestCasesByProblemID(int) ([]*TestCase, error)
	GetTestCaseSanityChecks(int) ([]*TestCase, error)
	GetTestCases() ([]*TestCase, error)
	UpdateTestCase(*TestCase) error

	// Submission CRU - no need for delete (yet)
	CreateSubmission(*Submission) (*Submission, error)
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
	dbHost := os.Getenv("DB_HOST")

	fmt.Println("Attemptng to connect to the server...")
	/* Create connection string */
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbUser, dbName, dbPass)

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
			encrypted_password VARCHAR(100),
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
            problem_name text,
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
			source_code TEXT,
			language INT,
			runtime_ms INT,
			mem_usage_kb INT
		)
	`

	_, err := s.db.Exec(query)
	return err
}

// -- Account Create --
func (s *PostgresStore) CreateAccount(acc *CreateAccountRequest) (*CreateAccountResponse, error) {
	query := `
			INSERT INTO Account (
                first_name,
                last_name,
				username,
				email,
				encrypted_password,
				created_at
			) 
			VALUES ($1, $2, $3, $4, $5, $6)
            RETURNING *
		`

	safePass, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Username, acc.Email, safePass, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	fmt.Println("Account inserted into database.")

	res, err := scanIntoAccountResponse(rows)

	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return res, nil
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

	return nil, fmt.Errorf("problem %d not found", id)
}

func (s *PostgresStore) GetProblemByName(name string) (*Problem, error) {
	query := `SELECT * FROM problem WHERE problem_name=$1`

	rows, err := s.db.Query(query, name)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoProblem(rows)
	}

	return nil, fmt.Errorf("Problem %s not found", name)
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
				io,
                is_sanity_check
			) 
			VALUES ($1, $2, $3, $4) RETURNING test_case_id;
		`

	var testCaseID int
	err := s.db.QueryRow(query, testcase.ProblemID, testcase.IO, testcase.IsSanityCheck).Scan(&testCaseID)
	if err != nil {
		return -1, err
	}

	return testCaseID, nil
}

// -- TestCase Read -- ID here is a PROBLEM id
func (s *PostgresStore) GetTestCasesByProblemID(id int) ([]*TestCase, error) {
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

// -- TestCase Read -- ID here is a PROBLEM id
func (s *PostgresStore) GetTestCaseSanityChecks(id int) ([]*TestCase, error) {
	query := `SELECT * FROM TestCase WHERE problem_id=$1 AND is_sanity_check=TRUE`

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
		fmt.Println("test case in storage", testCase)
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

// --  Submission Create --
func (s *PostgresStore) CreateSubmission(sub *Submission) (*Submission, error) {
	query := `
			INSERT INTO Submission (
				user_id,
				problem_id,
				submitted_at,
				code,
				language,
				exec_time_ms,
				mem_usage_kb
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

	rows, err := s.db.Query(query, sub.UserID, sub.ProblemID, time.Now().UTC(), sub.SourceCode, sub.Language)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	sub, err = scanIntoSubmission(rows)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// -- Submission Read --
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
	err := rows.Scan(&account.UserID, &account.FirstName, &account.LastName, &account.Username, &account.Email, &account.Password, &account.CreatedAt)

	return account, err
}

func scanIntoAccountResponse(rows *sql.Rows) (*CreateAccountResponse, error) {
	defer rows.Close()
	acc := new(Account)

	for rows.Next() {
		err := rows.Scan(&acc.UserID, &acc.FirstName, &acc.LastName, &acc.Username, &acc.Email, &acc.Password, &acc.CreatedAt)

		if err != nil {
			return nil, err
		}
	}

	res := &CreateAccountResponse{
		FirstName: acc.FirstName,
		LastName:  acc.LastName,
		Username:  acc.Username,
		Email:     acc.Email,
		Password:  acc.Password, // TODO: ENCRYPT THIS
		CreatedAt: acc.CreatedAt,
	}

	return res, nil
}

func scanIntoSubmission(rows *sql.Rows) (*Submission, error) {
	sub := new(Submission)
	err := rows.Scan(&sub.UserID, &sub.ProblemID, &sub.SubmittedAt, &sub.SourceCode, &sub.Language)

	return sub, err
}

func scanIntoTestCase(rows *sql.Rows) (*TestCase, error) {
	tc := new(TestCase)
	var ioData []byte

	err := rows.Scan(&tc.TestCaseID, &tc.ProblemID, &tc.IsSanityCheck, &ioData)
	if err != nil {
		return nil, err
	}

	fmt.Println(tc)

	err = json.Unmarshal(ioData, &tc.IO)
	if err != nil {
		return nil, err
	}

	return tc, nil
}

func scanIntoProblem(rows *sql.Rows) (*Problem, error) {
	p := new(Problem)
	err := rows.Scan(&p.ProblemID, &p.Prompt, &p.StarterCode, &p.Difficulty, &p.ProblemName, &p.FunctionName)

	return p, err
}
