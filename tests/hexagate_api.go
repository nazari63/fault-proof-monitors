package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	VALIDATE_ENDPOINT = "https://api.hexagate.com/api/v1/invariants/validate"
)

type ValidateRequest struct {
	Gate    string         `json:"gate"`
	ChainId int            `json:"chain_id"`
	Params  map[string]any `json:"params"`
	Mocks   map[string]any `json:"mocks"`
	Trace   bool           `json:"trace"`
}

type ValidateResponse struct {
	Count      int   `json:"count"`
	Failed     []any `json:"failed"`
	Exceptions []any `json:"exceptions"`
	Trace      any   `json:"trace"`
}

func ReadGateFile(filename string) (string, error) {
	file, err := os.Open(fmt.Sprintf("../monitors/%s", filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data[:]), nil
}

func HandleValidateRequest(gatefile string, params map[string]any, mocks map[string]any) ([]any, []any, any, error) {
	requestData := ValidateRequest{
		Gate:    gatefile,
		ChainId: 1,
		Params:  params,
		Mocks:   mocks,
		Trace:   true,
	}

	// marshal data into expected JSON format
	data := new(bytes.Buffer)
	enc := json.NewEncoder(data)
	enc.SetEscapeHTML(false)
	err := enc.Encode(requestData)
	if err != nil {
		return []any{}, []any{}, nil, err
	}

	// retrieve API key from .env file and call the API endpoint
	err = godotenv.Load("../.env")
	if err != nil {
		return []any{}, []any{}, nil, err
	}

	// create the POST request
	req, err := http.NewRequest("POST", VALIDATE_ENDPOINT, data)
	if err != nil {
		return []any{}, []any{}, nil, err
	}

	// set the appliable headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hexagate-Api-Key", os.Getenv("HEXAGATE_API_KEY"))

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []any{}, []any{}, nil, err
	}
	defer resp.Body.Close()

	// parse and decode the response
	var response ValidateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return []any{}, []any{}, nil, err
	}

	// access the response data
	failed := response.Failed
	exceptions := response.Exceptions
	trace := response.Trace

	return failed, exceptions, trace, nil
}
