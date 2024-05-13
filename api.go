package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	handler := cors.Default().Handler(router)
	router.Use(commonMiddleware)

	/* Run code */
	router.HandleFunc("/run", makeHTTPHandlerFunc(s.handleRunCode))

	/* Accounts */
	router.HandleFunc("/accounts", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/accounts/{id}", makeHTTPHandlerFunc(s.handleAccountByID))

	/* Problems */
	router.HandleFunc("/problems", makeHTTPHandlerFunc(s.handleProblem))
	router.HandleFunc("/problems/{id}", makeHTTPHandlerFunc(s.handleProblemByID))

	/* Test Cases */
	router.HandleFunc("/testcases", makeHTTPHandlerFunc(s.handleTestCase))
	router.HandleFunc("/testcases/{id}", makeHTTPHandlerFunc(s.handleTestCaseByProblemID))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, handler)
}

// Pass WriteJSON a pointer to a struct as param `v`, not sure what other types would work, if any ?
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request, idType string) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return id, fmt.Errorf("Invalid %s", idType)
	}

	return id, nil
}
