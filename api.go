package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ApiClient struct {
	Client  *http.Client
	Gateway string
	Token   string
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewApiClient(gateway string) *ApiClient {
	return &ApiClient{
		Client:  &http.Client{},
		Gateway: gateway,
	}
}

func (ac *ApiClient) Login(username, password string) (err error) {
	ac.Token = ""

	response, err := ac.DoRequest("signin", Credentials{username, password})
	if err != nil {
		err = fmt.Errorf("login request failed: %w", err)
		return
	}

	var token AuthTokenAnswer

	err = json.Unmarshal(response, &token)
	if err != nil {
		err = fmt.Errorf("could not parse token from login answer: %w", err)
		return
	}

	ac.Token = token.Token

	return
}

func (ac *ApiClient) DoRequest(url string, body interface{}) (respBody []byte, err error) {
	// Build request body as JSON
	var buf bytes.Buffer

	j := json.NewEncoder(&buf)

	err = j.Encode(body)
	if err != nil {
		err = fmt.Errorf("could not build JSON request data: %w", err)
		return
	}

	// Build Request
	baseUrl := "https://" + ac.Gateway + "/api/"

	req, err := http.NewRequest("POST", baseUrl+url, &buf)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	if ac.Token != "" {
		req.Header.Add("Authorization", "Bearer "+ac.Token)
	}

	resp, err := ac.Client.Do(req)
	if err != nil {
		err = fmt.Errorf("executing API request failed: %w", err)
		return
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("reading API response failed: %w", err)
		return
	}

	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with status %d", resp.StatusCode)
		return
	}

	return
}
