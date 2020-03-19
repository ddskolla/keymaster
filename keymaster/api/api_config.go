package api

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// api config

// roles
// idp
// creds
// client
// workflow endpoints
// env endpoints
// access control

// cli config

type ApiConfig struct {
	Name        string                `json:"name"`
	Idp         IdpConfig             `json:"idp"`
	Roles       map[string]RoleConfig `json:"roles"`
	Credentials []CredentialsConfig   `json:"credentials"`
	Workflow    WorkflowConfig        `json:"workflow"`
}

type IdpConfig struct {
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
}

func (c *IdpConfig) UnmarshalJSON(data []byte) error {
	var t struct {
		Type          string          `json:"type"`
		UntypedConfig json.RawMessage `json:"config"`
	}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	c.Type = t.Type
	switch c.Type {
	case "saml":
		var typedConfig IdpConfigSaml
		err := json.Unmarshal(t.UntypedConfig, &typedConfig)
		if err != nil {
			return err
		}
		c.Config = typedConfig
		return nil
	case "oidc":
		var typedConfig IdpConfigOidc
		err := json.Unmarshal(t.UntypedConfig, &typedConfig)
		if err != nil {
			return err
		}
		c.Config = typedConfig
		return nil
	}
	return errors.New("unknown idp type: " + c.Type)
}

type IdpConfigSaml struct {
	Issuer      string
	Audience    string
	Certificate string
}

type IdpConfigOidc struct {
	// Not implemented
}

type RoleConfig struct {
	Name string
}

type CredentialsConfig struct {
	Name string
	Type string
	Config interface{}
}

func (c *CredentialsConfig) UnmarshalJSON(data []byte) error {
	var t struct {
		Name string `json:"name"`
		Type          string          `json:"type"`
		UntypedConfig json.RawMessage `json:"config"`
	}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	c.Name = t.Name
	c.Type = t.Type
	switch c.Type {
	case "ssh_ca":
		var typedConfig CredentialsConfigSSH
		err := json.Unmarshal(t.UntypedConfig, &typedConfig)
		if err != nil {
			return err
		}
		c.Config = typedConfig
		return nil
	case "kubernetes":
		var typedConfig CredentialsConfigKube
		err := json.Unmarshal(t.UntypedConfig, &typedConfig)
		if err != nil {
			return err
		}
		c.Config = typedConfig
		return nil
	case "iam_assume_role":
		var typedConfig CredentialsConfigIAMAssumeRole
		err := json.Unmarshal(t.UntypedConfig, &typedConfig)
		if err != nil {
			return err
		}
		c.Config = typedConfig
		return nil
	case "iam_user":
		var typedConfig CredentialsConfigIAMUser
		err := json.Unmarshal(t.UntypedConfig, &typedConfig)
		if err != nil {
			return err
		}
		c.Config = typedConfig
		return nil
	}
	return errors.New("unknown credential type: " + c.Type)
}


type CredentialsConfigSSH struct {
	CAKey string
}

type CredentialsConfigKube struct {
	CAKey string
}

type CredentialsConfigIAMAssumeRole struct {
	// TODO
}

type CredentialsConfigIAMUser struct {
	// TODO
}

type WorkflowConfig struct {
	// TODO
}

type AccessControlConfig struct {
	WhiteListCidrs []string `json:"whitelist_cidrs"`
}
