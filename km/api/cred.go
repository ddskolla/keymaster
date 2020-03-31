package api

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Cred struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Expiry int64       `json:"expiry"`
	Value  interface{} `json:"value"`
}

type SSHCred struct {
	Username    string `json:"username"`
	Certificate []byte `json:"certficate"`
	PrivateKey  []byte `json:"private_key"`
}

type KubeCred struct {
	Username   string `json:"username"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type IAMCred struct {
	ProfileName     string `json:"profile_name"`
	RoleArn         string `json:"role_arn"`
	RoleSessionName string `json:"role_session_name"`
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
}

func (c *Cred) UnmarshalJSON(data []byte) error {
	var t struct {
		Name         string      `json:"name"`
		Type         string      `json:"type"`
		Expiry       int64       `json:"expiry"`
		UntypedValue json.RawMessage `json:"value"`
	}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	c.Name = t.Name
	c.Type = t.Type
	c.Expiry = t.Expiry
	var v interface{}
	switch c.Type {
	case "ssh":
		v = &SSHCred{}
	case "kube":
		v = &KubeCred{}
	case "iam":
		v = &IAMCred{}
	default:
		return errors.New("unknown credential type: " + c.Type)
	}
	err = json.Unmarshal(t.UntypedValue, v)
	if err != nil {
		return err
	}
	c.Value = v
	return nil
}
