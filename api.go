package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ApiClient struct {
	Client  *http.Client
	Gateway string
	Token   string
	Timeout time.Duration
	UseTls  bool
}

const DefaultTimeout = 5

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewApiClient(gateway string) *ApiClient {
	return &ApiClient{
		Client:  &http.Client{},
		Gateway: gateway,
		Timeout: DefaultTimeout * time.Second,
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

func (ac *ApiClient) DoLegacyRequest(mode string,
	to string,
	text string,
	username string,
	password string) error {
	params := url.Values{}

	switch mode {
	case "contactgroup":
		params.Add("mode", "group")
	case "contact":
		params.Add("mode", "user")
	default:
		params.Add("mode", "number")
	}

	params.Add("to", to)
	params.Add("text", text)
	params.Add("username", username)
	params.Add("password", password)

	schema := "http://"

	if ac.UseTls {
		schema = "https://"
	}

	myUrl := schema + ac.Gateway + "/api.php" + "?" + params.Encode()

	// Setup Timeout context
	ctx, cancel := context.WithTimeout(context.Background(), ac.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", myUrl, nil)

	if err != nil {
		return err
	}

	resp, err := ac.Client.Do(req)

	if err != nil {
		// We want to override the context error message to be more expressive
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("timeout during HTTP request: %w", err)
		}

		return fmt.Errorf("executing API request failed: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("reading API response failed: %w\nBody: %s", err, respBody)
		return err
	}

	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with status %d", resp.StatusCode)
		return err
	}

	return nil
}

func (ac *ApiClient) DoRequest(rawUrl string, body interface{}) (respBody []byte, err error) {
	// Build request body as JSON
	var buf bytes.Buffer

	j := json.NewEncoder(&buf)

	err = j.Encode(body)
	if err != nil {
		err = fmt.Errorf("could not build JSON request data: %w", err)
		return
	}

	// Setup Timeout context
	ctx, cancel := context.WithTimeout(context.Background(), ac.Timeout)
	defer cancel()

	schema := "http://"

	if ac.UseTls {
		schema = "https://"
	}

	// Build Request
	baseUrl := schema + ac.Gateway + "/api/"

	req, err := http.NewRequestWithContext(ctx, "POST", baseUrl+rawUrl, &buf)

	if err != nil {
		return []byte(""), err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	if ac.Token != "" {
		req.Header.Add("Authorization", "Bearer "+ac.Token)
	}

	resp, err := ac.Client.Do(req)

	if err != nil {
		// We want to override the context error message to be more expressive
		if errors.Is(err, context.DeadlineExceeded) {
			return []byte(""), fmt.Errorf("timeout during HTTP request: %w", err)
		}

		return []byte(""), fmt.Errorf("executing API request failed: %w", err)
	}

	respBody, err = io.ReadAll(resp.Body)
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
