package main

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestApiClient_Login(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://brevisone.local/api/signin",
		httpmock.NewStringResponder(200, `{"jwt":"abc123","expireAt":0}`))

	ac := NewApiClient("brevisone.local")

	err := ac.Login("admin", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	//nolint:all
	if "abc123" != ac.Token {
		t.Fatalf("expected %v, got %v", "abc123", ac.Token)
	}
}

func TestApiClient_LoginTimeout(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://brevisone.local/api/signin",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `{"jwt":"abc123","expireAt":0}`)

			time.Sleep(3 * time.Second)

			return resp, nil
		},
	)

	ac := NewApiClient("brevisone.local")

	ac.Timeout = 1 * time.Second
	err := ac.Login("admin", "password")
	// Validate that the error message is what we defined
	if !strings.Contains(err.Error(), "timeout during HTTP request") {
		t.Fatalf("expected %v, got %v", "timeout during HTTP request", err.Error())
	}
}

func TestApiClient_LoginErr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://brevisone.local/api/signin",
		httpmock.NewStringResponder(401, `{}`))

	ac := NewApiClient("brevisone.local")

	err := ac.Login("admin", "password")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestApiClient_UnmarshalErr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://brevisone.local/api/signin",
		httpmock.NewStringResponder(200, `{`))

	ac := NewApiClient("brevisone.local")

	err := ac.Login("admin", "password")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestApiClient_DoRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://brevisone.local/api/test",
		httpmock.NewStringResponder(200, `{"test":true}`))

	ac := NewApiClient("brevisone.local")
	ac.Token = "abc1234"

	response, err := ac.DoRequest("test", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	//nolint:all
	if `{"test":true}` != string(response) {
		t.Fatalf("expected %v, got %v", `{"test":true}`, string(response))
	}
}

func TestApiClient_DoLegacyRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://brevisone.local/api.php?mode=number&password=password&text=text&to=to&username=username",
		httpmock.NewStringResponder(200, `{"test":true}`))

	ac := NewApiClient("brevisone.local")
	ac.Token = "abc1234"

	err := ac.DoLegacyRequest("test", "to", "text", "username", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestApiClient_DoRequestErr(t *testing.T) {
	ac := NewApiClient("local")

	_, err := ac.DoRequest("test", nil)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}
