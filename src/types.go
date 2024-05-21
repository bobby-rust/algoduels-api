package main

import (
	"time"
)

type DifficultyRegistry struct {
	Easy   uint8
	Medium uint8
	Hard   uint8
}

func newDifficultyRegistry() *DifficultyRegistry {
	return &DifficultyRegistry{
		Easy:   1,
		Medium: 2,
		Hard:   3,
	}
}

type Account struct {
	UserID    int       `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"` // For some reason, need this json struct tag is needed to keep the formatting from giving me OCD...
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Problem struct {
	ProblemID   int    `json:"problem_id"`
	ProblemName string `json:"problem_name"`
	Prompt      string `json:"prompt"`
	StarterCode string `json:"starter_code"`
	Difficulty  uint8  `json:"difficulty"`
}

type TestCase struct {
	TestCaseID    int  `json:"test_case_id"`
	ProblemID     int  `json:"problem_id"`
	IO            IO   `json:"io"`
	IsSanityCheck bool `json:"is_sanity_check"`
}

type IO struct {
	Input  map[string]interface{} `json:"input"`
	Output interface{}            `json:"output"`
}

type Submission struct {
	SubmissionID int       `json:"submission_id"`
	UserID       int       `json:"user_id"`
	ProblemID    int       `json:"problem_id"`
	SubmittedAt  time.Time `json:"submitted_at"`
	SourceCode   string    `json:"source_code"`
	Language     int       `json:"language"`
	RuntimeMs    int       `json:"runtime_ms"`
	MemUsageKb   int       `json:"mem_usage_kb"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type CreateAccountResponse struct {
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateProblemRequest struct {
	ProblemName string `json:"problem_name"`
	Prompt      string
	StarterCode string `json:"starter_code"`
	Difficulty  int
}

type CreateTestCaseRequest struct {
	ProblemID     int  `json:"problem_id"`
	IO            IO   `json:"io"`
	IsSanityCheck bool `json:"is_sanity_check"`
}

type CreateSubmissionRequest struct {
	UserID     int    `json:"user_id"`
	ProblemID  int    `json:"problem_id"`
	SourceCode string `json:"source_code"`
	Language   int    `json:"language"`
	RuntimeMs  int32  `json:"runtime_ms"`
	MemUsageKb int32  `json:"mem_usage_kb"`
}

func NewAccountResponse(username, firstName, lastName, email, password string) *CreateAccountResponse {
	return &CreateAccountResponse{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().UTC(),
	}
}

func NewAccountRequest(username, firstName, lastName, email, password string) *CreateAccountRequest {
	return &CreateAccountRequest{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}

func NewProblem(problemName, prompt, starterCode string, difficulty uint8) *Problem {
	return &Problem{
		ProblemName: problemName,
		Prompt:      prompt,
		StarterCode: starterCode,
		Difficulty:  difficulty,
	}
}

func NewTestCase(problemID int, input map[string]interface{}, output interface{}, isSanityCheck bool) *TestCase {
	return &TestCase{
		ProblemID:     problemID,
		IO:            IO{Input: input, Output: output},
		IsSanityCheck: isSanityCheck,
	}
}

func NewSubmission(userID, problemID int, token, code string, language int, isAccepted bool, execTimeMs int, memUsageKb int) *Submission {
	return &Submission{
		UserID:     userID,
		ProblemID:  problemID,
		SourceCode: code,
		Language:   language,
		RuntimeMs:  execTimeMs,
		MemUsageKb: memUsageKb,
	}
}
