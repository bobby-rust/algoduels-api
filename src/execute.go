package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	url       = "http://localhost:2358/submissions" // judge0 url
	urlParams = "?&fields=stdout,time,memory,stderr,compile_output,message,status"
)

var languageIDs = map[string]int{
	"python3":    71,
	"javascript": 63,
}

type ExecReq struct {
	ProblemID     int    `json:"problem_id"`
	LanguageID    int    `json:"language_id"`
	SourceCode    string `json:"source_code"`
	IsSanityCheck bool   `json:"is_sanity_check"`
}

type Result struct {
	Passed bool         `json:"passed"`
	Result []TestResult `json:"result"`
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
	res, err := http.Post(url+urlParams, "application/json", bytes.NewReader(jsonReq)) // http.Post takes io.Reader for the request body
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
 * Gets all test cases for a problem id, calls test(), returns Result
 */
func submit(s *PostgresStore, req ExecReq) (*Result, error) {
	result := new(Result)

	return result, nil
}

/**
 * Gets sanity check test cases for a problem id, calls test(), returns Result
 */
func run(s *PostgresStore, req ExecReq) (*Result, error) {
	result := new(Result)

	return result, nil
}

/**
 *
 */
func test(testCases []TestCase) *[]Result {
	results := new([]Result)

	return results
}

/* Polls judge0 to retreive the results of the submission associated with `token` */
func pollJudge0Submission(token string) (*ExecResult, error) {
	timeout := 10 * time.Second
	startTime := time.Now()
	time.Sleep(time.Second)
	ran := 0

	for time.Since(startTime) < timeout {
		fmt.Println("Giving submission time to process...")
		time.Sleep(time.Second * 1)
		fmt.Println("Sending GET request...")
		outputResp, err := http.Get(url + "/" + token)
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
			fmt.Println("Status.Description: ", outputRespStruct.Status.Description)
			return nil, errors.New("Error description: " + outputRespStruct.Status.Description)
		}
	}

	return nil, errors.New("Time limit exceeded")
}
