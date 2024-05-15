package main

import (
	"fmt"
	"net/http"
)

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch method := r.Method; method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		return s.handleGetAccountByID(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleProblem(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		return s.handleGetProblems(w, r)
	case "POST":
		return s.handleCreateProblem(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleProblemByID(w http.ResponseWriter, r *http.Request) error {
	if method := r.Method; method == "GET" {
		return s.handleGetProblemByID(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleTestCase(w http.ResponseWriter, r *http.Request) error {
	if method := r.Method; method == "POST" {
		return s.handleCreateTestCase(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleTestCaseByProblemID(w http.ResponseWriter, r *http.Request) error {
	if method := r.Method; method == "GET" {
		return s.handleGetTestCasesByProblemID(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleSubmission(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		return s.handleGetSubmissions(w, r)
	case "POST":
		return s.handleCreateSubmission(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleSubmissionByID(w http.ResponseWriter, r *http.Request) error {
	if method := r.Method; method == "GET" {
		return s.handleGetSubmissionByID(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleRunCode(w http.ResponseWriter, r *http.Request) error {
	if method := r.Method; method == "POST" {
		return execute(w, r)
	}

	return fmt.Errorf("Method not supported %s", r.Method)
}
