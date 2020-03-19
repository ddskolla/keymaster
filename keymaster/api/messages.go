package api

type ConfigRequest struct {
}

type ConfigResponse struct {
	Version string `json:"version"`
}

type DirectSamlAuthRequest struct {
	RequestedAccess string  `json:"requested_access"`
	SAMLResponse    string  `json:"saml_response"`
	SigAlg          string  `json:"sig_alg"`
	Signature       string  `json:"signature"`
	RelayState      *string `json:"relay_state,omitempty"`
}

type DirectOidcAuthRequest struct {
}

type DirectAuthResponse struct {
	Credentials map[string][]byte `json:"result"`
}

type WorkflowStartRequest struct {
}

type WorkflowStartResponse struct {
}

type WorkflowAuthRequest struct {
}

type WorkflowAuthResponse struct {
	Credentials map[string][]byte `json:"result"`
}
