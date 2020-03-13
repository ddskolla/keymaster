package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestUnmarshalJSON(t *testing.T) {
	var req Request
	var err error

	// Invalid json blech
	for _, payload := range []string{"", "(╯°□°)╯︵ ┻━┻", "}", "whatever"} {
		err = json.Unmarshal([]byte(payload), &req)
		if assert.Error(t, err) {
			assert.NotEmpty(t, err.Error())
		}
	}

	// Missing or bad op
	for _, payload := range []string{"{}", "{ \"op\": \"Unknown\" }"} {
		err = json.Unmarshal([]byte(payload), &req)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "request: unknown operation")
		}
	}

	// All the supposed to be valid things
	var validCases = []struct {
		Input          string
		ExpectedResult interface{}
	}{
		{Input: `{ "op": "ping", "message": {} }`,
			ExpectedResult: Request{
				Operation: "ping",
				Message: PingRequestMessage{},
			},
		},
		{Input: `{ "op": "direct_saml_auth", "message": { "requested_access": "foo" } }`,
			ExpectedResult: Request{
				Operation: "direct_saml_auth",
				Message:   DirectSamlAuthRequestMessage{RequestedAccess:"foo"},
			},
		},
	}
	for _, testCase := range validCases {
		var req Request
		err = json.Unmarshal([]byte(testCase.Input), &req)
		if assert.NoError(t, err) {
			assert.Equal(t, testCase.ExpectedResult, req)
		}
	}
}
