package api

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Request struct {
	Operation string      `json:"op"`
	Message   interface{} `json:"message"`
}

type PingRequestMessage struct {
}

type PingResponseMessage struct {
	Message string
}

type DirectSamlAuthRequestMessage struct {
	RequestedAccess string  `json:"requested_access"`
	SAMLResponse    string  `json:"saml_response"`
	SigAlg          string  `json:"sig_alg"`
	Signature       string  `json:"signature"`
	RelayState      *string `json:"relay_state,omitempty"`
}

type DirectOidcAuthRequestMessage struct {
	// TODO
}

type DirectAuthResponseMessage struct {
	Credentials map[string][]byte `json:"result"`
}

func (r *Request) UnmarshalJSON(b []byte) error {
	var tmp struct {
		Operation string          `json:"op"`
		Message   json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	r.Operation = tmp.Operation
	switch r.Operation {
	case "ping":
		var m PingRequestMessage
		if err := json.Unmarshal(tmp.Message, &m); err != nil {
			return errors.Wrap(err, "request: invalid ping message")
		}
		r.Message = m
		return nil
	case "direct_saml_auth":
		var m DirectSamlAuthRequestMessage
		if err := json.Unmarshal(tmp.Message, &m); err != nil {
			return errors.Wrap(err, "request: invalid direct_saml_auth message")
		}
		r.Message = m
		return nil
	}
	return errors.New("request: unknown operation")
}
