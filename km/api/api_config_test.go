package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRequest_UnmarshalJSON(t *testing.T) {
	testCases := map[string]Request{
		"config": {
			Type: "config",
			Payload: &ConfigRequest{},
		},
	}

	// Unmarshal c -> c2, check c == c2
	for _, c := range testCases {
		b, err := json.Marshal(c)
		log.Println(string(b))
		assert.NoError(t, err)
		assert.NotEmpty(t, b)

		var c2 Request
		err = json.Unmarshal(b, &c2)
		assert.NoError(t, err)

		assert.Equal(t, c, c2)
	}
}

func TestIdpConfig_UnmarshalJSON(t *testing.T) {
	testCases := map[string]IdpConfig{
		"t1": {
			Type: "saml",
			Config: &IdpConfigSaml{
				Issuer: "foo",
				Audience: "bar",
				Certificate: "pem-goes-here",
			},
		},
	}

	// Unmarshal c -> c2, check c == c2
	for _, c := range testCases {
		b, err := json.Marshal(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, b)

		var c2 IdpConfig
		err = json.Unmarshal(b, &c2)
		assert.NoError(t, err)

		assert.Equal(t, c, c2)
	}
}

func TestCredentialsConfig_UnmarshalJSON(t *testing.T) {
	testCases := map[string]CredentialsConfig{
		"ssh1": {
			Name: "ssh-example",
			Type: "ssh_ca",
			Config: &CredentialsConfigSSH{
				CAKey:"my-ssh-ca-key",
			},
		},
		"kube1": {
			Name: "kube-example",
			Type: "kubernetes",
			Config: &CredentialsConfigKube{
			},
		},
		"iam_assumerole1": {
			Name: "iam-assumerole-example",
			Type: "iam_assume_role",
			Config: &CredentialsConfigIAMAssumeRole{
			},
		},
		"iam_user1": {
			Name: "iam-user-example",
			Type: "iam_user",
			Config: &CredentialsConfigIAMUser{
			},
		},
	}

	// Unmarshal c -> c2, check c == c2
	for _, c := range testCases {
		b, err := json.Marshal(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, b)

		var c2 CredentialsConfig
		err = json.Unmarshal(b, &c2)
		assert.NoError(t, err)

		assert.Equal(t, c, c2)
	}
}