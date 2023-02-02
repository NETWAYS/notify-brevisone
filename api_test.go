package main

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestApiClient_Login(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://brevisone.local/api/signin",
		httpmock.NewStringResponder(200, `{"jwt":"abc123","expireAt":0}`))

	ac := NewApiClient("brevisone.local")

	err := ac.Login("admin", "password")
	assert.NoError(t, err)

	assert.Equal(t, "abc123", ac.Token)
}

func TestApiClient_DoRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://brevisone.local/api/test",
		httpmock.NewStringResponder(200, `{"test":true}`))

	ac := NewApiClient("brevisone.local")
	ac.Token = "abc1234"

	response, err := ac.DoRequest("test", nil)
	assert.NoError(t, err)
	assert.Equal(t, `{"test":true}`, string(response))
}
