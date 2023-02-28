package main

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiClient_Login(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://brevisone.local/api/signin",
		httpmock.NewStringResponder(200, `{"jwt":"abc123","expireAt":0}`))

	ac := NewApiClient("brevisone.local")

	err := ac.Login("admin", "password")
	assert.NoError(t, err)

	assert.Equal(t, "abc123", ac.Token)
}

func TestApiClient_LoginErr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://brevisone.local/api/signin",
		httpmock.NewStringResponder(401, `{}`))

	ac := NewApiClient("brevisone.local")

	err := ac.Login("admin", "password")
	assert.Error(t, err)
}

func TestApiClient_UnmarshalErr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://brevisone.local/api/signin",
		httpmock.NewStringResponder(200, `{`))

	ac := NewApiClient("brevisone.local")

	err := ac.Login("admin", "password")
	assert.Error(t, err)
}

func TestApiClient_DoRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://brevisone.local/api/test",
		httpmock.NewStringResponder(200, `{"test":true}`))

	ac := NewApiClient("brevisone.local")
	ac.Token = "abc1234"

	response, err := ac.DoRequest("test", nil)
	assert.NoError(t, err)
	assert.Equal(t, `{"test":true}`, string(response))
}

func TestApiClient_DoRequestErr(t *testing.T) {
	ac := NewApiClient("local")

	_, err := ac.DoRequest("test", nil)
	assert.Error(t, err)
}
