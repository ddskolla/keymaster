package api

import (
	"encoding/json"
	"github.com/bsycorp/keymaster/km/util"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"strings"
)

type Server struct {
	Config Config
}

func (s *Server) Configure(config string) error {
	var err error
	var tmpConfig Config

	// Load config, maybe "by reference" (config env var might be a
	// literal or a reference to a bucket or file).
	// Then we load as YAML or JSON.
	configData, err := util.Load(config)
	if strings.HasPrefix(string(configData), "{") {
		err = json.Unmarshal(configData, &tmpConfig)
		if err != nil {
			return err
		}
	} else {
		err = yaml.Unmarshal(configData, &tmpConfig)
		if err != nil {
			return err
		}
	}
	s.Config = tmpConfig
	return nil
}

func (s *Server) HandleConfig(req *ConfigRequest) (*ConfigResponse, error) {
	// Copy the public parts of our configuration.
	var resp ConfigResponse
	resp.Config = ConfigPublic{
		Name:     s.Config.Name,
		Idp:      s.Config.Idp,
		Roles:    s.Config.Roles,
		Workflow: s.Config.Workflow,
	}
	return &resp, nil
}

func (s *Server) HandleDirectSamlAuth(eq *DirectSamlAuthRequest) (*DirectAuthResponse, error) {
	return nil, errors.New("Not implemented")
}

func (s *Server) HandleDirectOidcAuth(req *DirectOidcAuthRequest) (*DirectAuthResponse, error) {
	return nil, errors.New("Not implemented")
}

func (s *Server) HandleWorkflowStart(req *WorkflowStartRequest) (*WorkflowStartResponse, error) {
	return nil, errors.New("Not implemented")
}

func (s *Server) HandleWorkflowAuth(req *WorkflowAuthRequest) (*WorkflowAuthResponse, error) {
	return nil, errors.New("Not implemented")
}
