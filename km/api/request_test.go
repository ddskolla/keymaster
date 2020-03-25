package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest_UnmarshalJSON(t *testing.T) {
	testCases := map[string]Request{
		"config": {
			Type: "config",
			Payload: &ConfigRequest{},
		},
		"direct_saml_auth": {
			Type: "direct_saml_auth",
			Payload: &DirectSamlAuthRequest{},
		},
		"direct_oidc_auth": {
			Type: "direct_oidc_auth",
			Payload: &DirectOidcAuthRequest{},
		},
		"workflow_start": {
			Type: "workflow_start",
			Payload: &WorkflowStartRequest{},
		},
		"workflow_auth": {
			Type: "workflow_auth",
			Payload: &WorkflowAuthRequest{},
		},
	}

	// Unmarshal c -> c2, check c == c2
	for _, c := range testCases {
		b, err := json.Marshal(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, b)

		var c2 Request
		err = json.Unmarshal(b, &c2)
		assert.NoError(t, err)

		assert.Equal(t, c, c2)
	}
}