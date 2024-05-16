package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GET api/users/{id}
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r, "user_id")
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// GET api/users
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// POST api/users
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var acc CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		return err
	}
	defer r.Body.Close()

	accountReq := NewAccountRequest(acc.Username, acc.FirstName, acc.LastName, acc.Email, acc.Password)
	accountRes, err := s.store.CreateAccount(accountReq)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, accountRes)
}

// DELETE api/users
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r, "user_id")
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusNoContent, map[string]int{"deleted": id})
}

// GET api/problems/{id}
func (s *APIServer) handleGetProblemByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r, "problem_id")
	if err != nil {
		return err
	}

	p, err := s.store.GetProblemByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, p)
}

// GET api/problems
func (s *APIServer) handleGetProblems(w http.ResponseWriter, r *http.Request) error {
	problems, err := s.store.GetProblems()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, problems)
}

// POST api/problems
func (s *APIServer) handleCreateProblem(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateProblemRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	defer r.Body.Close()
	fmt.Printf("in handler: prompt %s, starterCode %s, difficulty %d\n", req.Prompt, req.StarterCode, req.Difficulty)
	problem := NewProblem(req.Prompt, req.StarterCode, uint8(req.Difficulty))
	problemID, err := s.store.CreateProblem(problem)
	if err != nil {
		return err
	}

	problem.ProblemID = problemID

	return WriteJSON(w, http.StatusCreated, problem)
}

// GET api/testcases/{id} *** id here is a PROBLEM id ***
func (s *APIServer) handleGetTestCasesByProblemID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r, "problem_id")
	if err != nil {
		return err
	}

	testCase, err := s.store.GetTestCaseByProblemID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, testCase)
}

// POST api/testcases
func (s *APIServer) handleCreateTestCase(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateTestCaseRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	defer r.Body.Close()

	testCase := NewTestCase(req.ProblemID, req.Input, req.Output, req.IsSanityCheck)

	id, err := s.store.CreateTestCase(testCase)
	if err != nil {
		return err
	}
	testCase.TestCaseID = id
	return WriteJSON(w, http.StatusCreated, testCase)
}

func (s *APIServer) handleCreateSubmission(w http.ResponseWriter, r *http.Request) error {
	req := new(Submission)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	sub, err := s.store.CreateSubmission(req)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, sub)
}

func (s *APIServer) handleGetSubmissionByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r, "submission_id")
	if err != nil {
		return err
	}

	testCase, err := s.store.GetSubmissionByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, testCase)
}

func (s *APIServer) handleGetSubmissions(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// POST api/submit
func (s *APIServer) handleSubmitCode(w http.ResponseWriter, r *http.Request) error {
	req := new(ExecReq)
}

// POST api/run
func (s *APIServer) handleRunCode(w http.ResponseWriter, r *http.Request) error {
	req := new(ExecReq)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	result, err := execute(req)
	if err != nil {
		return err
	}

	return nil
}
