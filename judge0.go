package main

import (
	"fmt"
	// "strconv"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

var languageIDs = map[string]int{
	"python3":    71,
	"javascript": 63,
}

type RunCodeRequest struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
}

type GetOutputReq struct {
	Token string `json:"token"`
}

type GetOutputResp struct {
	Stdout string `json:"stdout"`
	Time string `json:"time"`
	Memory string `json:"memory"`
	Stderr *string `json:"stderr"`
	Token string `json:"token"`
	CompileOutput string `json:"compile_output"`
	Message string `json:"message"`
	Status interface{} `json:"status"`

}
type CreateSubmissionResp struct {
	Token string `json:"token"`
}

func runCode(w http.ResponseWriter, r *http.Request) error {
	req := new(RunCodeRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	defer r.Body.Close()

	url := "http://localhost:2358/submissions"
	data := map[string]interface{} {
		"language_id": req.LanguageID,
		"source_code": req.SourceCode,
	}
	// reqBody := fmt.Sprintf(`{ "language_id": %s, "source_code": %s }`, req.LanguageID, req.SourceCode)
	reqBody, err := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response code: %d", resp.StatusCode)
	}

	// resBody := new(CreateSubmissionResp)
	// if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
	// 	return err
	// }
	//
	// fmt.Printf("%s", resBody)
	// Read the raw response body and print it as a string
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var createSubmissionResp CreateSubmissionResp
	if err := json.Unmarshal(responseBytes, &createSubmissionResp); err != nil {
		return err
	}
	responseString := string(responseBytes)
	fmt.Printf("Response Body:\n%s\n", responseString)
	
	token := createSubmissionResp.Token
	fmt.Println(createSubmissionResp)
	// getOutputUrlInt, err := fmt.Printf(url + "/%s", token)
	// if err != nil {
	// 	return err
	// }
	// getOutputUrl := strconv.Itoa(getOutputUrlInt)
	getOutputUrl := url + "/" + token
	fmt.Println(getOutputUrl)
	outputResp, err := http.Get(getOutputUrl)
	if err != nil {
		fmt.Println("Error during GET request: ")
		return err
	}
	getOutputRespBytes, err := ioutil.ReadAll(outputResp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(getOutputRespBytes))
	getOutputResp := new(GetOutputResp)
	outputErr := json.NewDecoder(outputResp.Body).Decode(&getOutputResp)
	if outputErr != nil {
		return err
	}
	fmt.Print("getOutputResp:")
	fmt.Println(*getOutputResp)
	return nil	 
}

func submitCode(w http.ResponseWriter, r *http.Request) error {
	return nil
}
