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
	ProblemID  int    `json:"problem_id"`
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	IsTest     bool   `json:"is_test"`
}

type Result struct {
	Passed bool         `json:"passed"`
	Result []TestResult `json:"result"`
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

/* Polls judge0 to retreive the results of the submission associated with `token`*/
func pollJudge0Submission(token string) (*ExecResult, error) {
	timeout := 10 * time.Second
	startTime := time.Now()
	time.Sleep(time.Second)
	ran := 0

	for time.Since(startTime) < timeout {
		fmt.Println("Giving submission time to process...")
		time.Sleep(time.Second * 1)
		fmt.Println("Sending GET request...")
		outputResp, err := http.Get(url)
		ran++
		if err != nil {
			fmt.Println("Error during GET request, retrying... ")
			continue
		}

		/* Parse get submission response */
		outputRespStruct := new(GetOutputResp)
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

/* Executes some code using Judge0 */
func execute(req *ExecReq) (*ExecResult, error) {

	j0CreateSubmissionReqBody := incomingReqBody                               // for ease of understanding, we copy the incoming req into a variable with a better name
	mJ0CreateSubmissionReqBody, err := json.Marshal(j0CreateSubmissionReqBody) // marshalled judge0 req body, we convert to raw byte slice for sending
	if err != nil {
		return err
	}

	/* Create judge0 code submission */
	j0CreateSubmissionResp, err := http.Post(url+urlParams, "application/json", bytes.NewReader(mJ0CreateSubmissionReqBody)) // http.Post takes io.Reader for the request body
	if err != nil {
		return err
	}
	defer j0CreateSubmissionResp.Body.Close()

	/* Parse judge0 create submission response */
	var j0CreateSubmissionRespStruct CreateSubmissionResp
	err = json.NewDecoder(j0CreateSubmissionResp.Body).Decode(&j0CreateSubmissionRespStruct)
	if err != nil {
		return err
	}

	/* Extract token */
	token := j0CreateSubmissionRespStruct.Token
	fmt.Println(token)

	/* url to retrieve code output */
	outputUrl := url + "/" + token
	fmt.Println(outputUrl)

	/* Poll judge0 until code has finished executing and output is ready */
	outputResp, err := pollJudge0Submission(outputUrl)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, outputResp)
}

func submitCode(w http.ResponseWriter, r *http.Request) error {
	return nil
}
