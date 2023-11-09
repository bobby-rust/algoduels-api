package main

import (
	"encoding/json"
	"net/http"
)

// GET api/users/{id}
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
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
	req := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	defer r.Body.Close()

	// account := &Account{} // same exact thing as new()
	account := NewAccount(req.Username, req.Email, req.Password)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}

// DELETE api/users
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
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
	id, err := getID(r)
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

	problem := NewProblem(req.Prompt, req.StarterCode, req.Difficulty)
	if err := s.store.CreateProblem(problem); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, problem)
}

// GET api/testcases/{id} *** id here is a PROBLEM id ***
func (s *APIServer) handleGetTestCaseByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	testCase, err := s.store.GetTestCaseByID(id)
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

	testCase := NewTestCase(req.ProblemID, req.Input, req.Output)
	if err := s.store.CreateTestCase(testCase); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, testCase)
}

func (s *APIServer) handleGetSubmissionByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	testCase, err := s.store.GetSubmissionByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, testCase)
}

func (s *APIServer) handleCreateSubmission(w http.ResponseWriter, r *http.Request) error {
	// req := new(CreateSubmissionRequest)
	//
	// if err := json.NewDecoder(r.Body).Decode(req); err != nil {
	// 	return err
	// }
	// defer r.Body.Close()
	//
	// // account := &Account{} // same exact thing as new()
	// sub := NewSubmission(req.UserID, req.ProblemID, req.Code, req.Language)
	// if err := s.store.CreateSubmission(sub); err != nil {
	// 	return err
	// }
	//
	// return WriteJSON(w, http.StatusCreated, sub)
	return nil
}

func (s* APIServer) handleGetSubmissions(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleSubmitCode(w http.ResponseWriter, r *http.Request) error {

	return nil
}
