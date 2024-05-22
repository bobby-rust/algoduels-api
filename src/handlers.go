package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ExecBatchReq struct {
	Submissions []ExecReq `json:"submissions"`
}

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

func (s *APIServer) handleGetProblemByName(w http.ResponseWriter, r *http.Request) error {
	name := mux.Vars(r)["name"]
	fmt.Println(name)
	problem, err := s.store.GetProblemByName(name)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, problem)
}

// POST api/problems
func (s *APIServer) handleCreateProblem(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateProblemRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	defer r.Body.Close()
	problem := NewProblem(req.ProblemName, req.Prompt, req.StarterCode, req.FunctionName, uint8(req.Difficulty))
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

	testCase, err := s.store.GetTestCasesByProblemID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, testCase)
}

func (s *APIServer) handleGetTestCaseSanityChecks(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r, "problem_id")
	if err != nil {
		return err
	}

	testCase, err := s.store.GetTestCaseSanityChecks(id)
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

	testCase := NewTestCase(req.ProblemID, req.IO.Input, req.IO.Output, req.IsSanityCheck)

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
	// req := new(ExecReq)

	return nil
}

// POST api/run
func (s *APIServer) handleRunCode(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("handling run code request...")
	req := new(ExecReq)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	result, err := run(s, req)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, result)
}

// POST api/run/batch
func (s *APIServer) handleRunBatchCode(w http.ResponseWriter, r *http.Request) error {
	req := new(ExecBatchReq)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("ErroR!!")
		return err
	}

	fmt.Println(req)

	var res []*ExecResult

	for i := 0; i < len(req.Submissions); i++ {
		result, err := run(s, &req.Submissions[i])
		if err != nil {
			return err
		}
		res = append(res, result)
	}

	return WriteJSON(w, http.StatusOK, res)
}
