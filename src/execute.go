package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// TODO: Return full judge0 response from API when an error is thrown

const (
	judge0Url       = "http://localhost:2358/submissions" // judge0 url
	judge0UrlParams = "?&fields=stdout,time,memory,stderr,compile_output,message,status"
	apiUrl          = "http://localhost:4000/api"
)

var languageIDs = map[string]int{
	"python3":    71,
	"javascript": 63,
}

type RunRes struct {
	Result  Result     `json:"result"`
	ExecRes ExecResult `json:"exec_res"`
}

type ExecReq struct {
	ProblemID     int    `json:"problem_id"`
	LanguageID    int    `json:"language_id"`
	SourceCode    string `json:"source_code"`
	IsSanityCheck bool   `json:"is_sanity_check"`
}

type Result struct {
	Passed      bool         `json:"passed"`
	TestResults []TestResult `json:"result"`
}

type CrSubRes struct {
	Token string `json:"token"`
}

type TestResult struct {
	Input    string `json:"input"`
	Output   string `json:"output"`
	Expected string `json:"expected"`
	Passed   bool   `json:"passed"`
}

type ExecResult struct {
	Stdout        string  `json:"stdout"`
	Time          string  `json:"time"`
	Memory        int     `json:"memory"`
	Stderr        *string `json:"stderr"`
	CompileOutput string  `json:"compile_output"`
	Message       string  `json:"message"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

/* Executes some code and returns result of execution */
func execute(req *ExecReq) (*ExecResult, error) {
	fmt.Println("executing...")
	jsonReq, err := json.Marshal(req) // marshalled (JSONified) judge0 req body, we convert to raw byte slice for sending
	if err != nil {
		fmt.Println("error marshalling json")
		return nil, err
	}

	fmt.Println("successfully marshalled json: ", jsonReq)

	/* Create judge0 code submission */
	res, err := http.Post(judge0Url+judge0UrlParams, "application/json", bytes.NewReader(jsonReq)) // http.Post takes io.Reader for the request body
	if err != nil {
		fmt.Println("Error sending code to judge0")
		return nil, err
	}
	defer res.Body.Close()

	/* Parse judge0 create submission response */
	var crSubRes CrSubRes
	err = json.NewDecoder(res.Body).Decode(&crSubRes)
	if err != nil {
		return nil, err
	}

	/* Extract token */
	token := crSubRes.Token
	fmt.Println(token)

	/* Poll judge0 until code has finished executing and output is ready */
	execResult, err := pollJudge0Submission(token)
	if err != nil {
		return nil, err
	}

	return execResult, nil
}

/**
 * Gets sanity check test cases for a problem id, calls test(), returns Result
 */
func run(s *APIServer, req *ExecReq) (*ExecResult, error) {
	result := new(ExecResult)

	fmt.Println(req.ProblemID)
	if req.IsSanityCheck {
		// fetch sanity tests
		tests, err := s.store.GetTestCaseSanityChecks(req.ProblemID)
		if err != nil {
			fmt.Println("ERROR!")
			return nil, err
		}

		fmt.Println(tests)
	} else {
		// fetch all tests
		tests, err := s.store.GetTestCasesByProblemID(req.ProblemID)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(tests); i++ {
			var inputValues string
			for _, val := range tests[i].IO.Input {
				inputValues += fmt.Sprintf("%v,", val)
			}
			inputValues = inputValues[:len(inputValues)-1] // remove last comma
			fmt.Println(tests[i])
		}
	}

	fmt.Println(req)
	fmt.Println(req.SourceCode)

	return result, nil
}

/**
 * executes code and checks if the results are correct
 */
func test(execReq *ExecReq, testCases *[]TestCase) *TestResult {

	result := new(TestResult)
	return result
}

/* Polls judge0 to retreive the results of the submission associated with `token` */
func pollJudge0Submission(token string) (*ExecResult, error) {
	timeout := 20 * time.Second
	startTime := time.Now()
	time.Sleep(time.Second)
	ran := 0

	for time.Since(startTime) < timeout {
		fmt.Println("Giving submission time to process...")
		time.Sleep(time.Second * 1)
		fmt.Println("Sending GET request...")
		outputResp, err := http.Get(judge0Url + "/" + token)
		ran++
		if err != nil {
			fmt.Println("Error during GET request, retrying... ")
			continue
		}

		/* Parse get submission response */
		outputRespStruct := new(ExecResult)
		outputErr := json.NewDecoder(outputResp.Body).Decode(&outputRespStruct)
		if outputErr != nil {
			fmt.Println("Error decoding outputResp body", outputErr)
			fmt.Printf("ran %d times", ran)
			return nil, outputErr
		}
		defer outputResp.Body.Close()

		if outputRespStruct.Status.ID == 3 {
			fmt.Println("Response: ", outputRespStruct)
			return outputRespStruct, nil
		} else {
			if outputRespStruct.Status.ID == 2 {
				continue
			}

			fmt.Println(outputRespStruct)
			fmt.Println("status.id: ", outputRespStruct.Status.ID)
			fmt.Println("Status.Description: ", outputRespStruct.Status.Description)
			return nil, errors.New("Error description: " + outputRespStruct.Status.Description)
		}
	}

	return nil, errors.New("Time limit exceeded")
}
