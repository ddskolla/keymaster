package api

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type ApiConfig struct {
	Name        string                `json:"name"`
	Idp         IdpConfig             `json:"idp"`
	Roles       map[string]RoleConfig `json:"roles"`
	Credentials []CredentialsConfig   `json:"credentials"`
	Workflow    WorkflowConfig        `json:"workflow"`
}

type IdpConfig struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
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
	Name            string `json:"name"`
	ValidForSeconds string `json:"valid_for_seconds"`
}

type CredentialsConfig struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
}

type CredentialsConfigSSH struct {
	CAKey      string   `json:"ca_key"`
	Principals []string `json:"principals"`
}

type CredentialsConfigKube struct {
	CAKey string `json:"ca_key"`
}

type CredentialsConfigIAMAssumeRole struct {
	TargetRole string `json:"target_role"`
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
	var config interface{}
	switch c.Type {
	case "saml":
		config = &IdpConfigSaml{}
	case "oidc":
		config = &IdpConfigOidc{}
	default:
		return errors.New("unknown idp type: " + c.Type)
	}
	err = json.Unmarshal(t.UntypedConfig, config)
	if err != nil {
		return err
	}
	c.Config = config
	return nil
}

func (c *CredentialsConfig) UnmarshalJSON(data []byte) error {
	var t struct {
		Name          string          `json:"name"`
		Type          string          `json:"type"`
		UntypedConfig json.RawMessage `json:"config"`
	}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	c.Name = t.Name
	c.Type = t.Type
	var config interface{}
	switch c.Type {
	case "ssh_ca":
		config = &CredentialsConfigSSH{}
	case "kubernetes":
		config = &CredentialsConfigKube{}
	case "iam_assume_role":
		config = &CredentialsConfigIAMAssumeRole{}
	case "iam_user":
		config = &CredentialsConfigIAMUser{}
	default:
		return errors.New("unknown credential type: " + c.Type)
	}
	err = json.Unmarshal(t.UntypedConfig, config)
	if err != nil {
		return err
	}
	c.Config = config
	return nil
}
