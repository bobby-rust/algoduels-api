package main

import (
	"time"
)

type Account struct {
	UserID        		int   	  `json:"user_id"`
	Username  			string 	  `json:"username"`
	Email	  			string 	  `json:"email"`
	EncryptedPassword	string	  `json:"encryptedPassword"`
	CreatedAt 			time.Time `json:"createdAt"`
}

type Problem struct {
	ProblemID 	int    `json:"problem_id"`
	Prompt		string `json:"prompt"`
	StarterCode string `json:"starter_code"`
	Difficulty 	int    `json:"difficulty"`
}

type TestCase struct {
	TestCaseID int 	  `json:"test_case_id"`
	ProblemID  int 	  `json:"problem_id"`
	Input 	   string `json:"input"`
	Output 	   string `json:"output"`
}

type Submission struct {
	SubmissionID int 	   `json:"submission_id"`
	UserID		 int 	   `json:"user_id"`
	ProblemID	 int 	   `json:"problem_id"`
	SubmittedAt	 time.Time `json:"submitted_at"`
	Code 		 string    `json:"code"`
	Language 	 int 	   `json:"language"`
	IsAccepted	 bool 	   `json:"is_accepted"`
	ExecTimeMS 	 int 	   `json:"exec_time_ms"`
	MemUsageKB 	 int 	   `json:"mem_usage_kb"`
}

type CreateAccountRequest struct {
	Username  string
	Email 	  string
	Password  string
}

func NewAccount(username, email, password string) *Account {
	return &Account{
		Username:  		   username,
		Email: 	   		   email,
		EncryptedPassword: password,
		CreatedAt: 		   time.Now().UTC(),
	}
}

func NewProblem()