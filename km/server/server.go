package server

import (
	"encoding/json"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/bsycorp/keymaster/km/creds"
	"github.com/bsycorp/keymaster/km/util"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
)

type Server struct {
	Config api.Config
}

func (s *Server) Configure(config string) error {
	var err error
	var tmpConfig api.Config

	// Load config, maybe "by reference" (config env var might be a
	// literal or a reference to a bucket or file).
	// Then we load as YAML or JSON.
	configData, err := util.Load(config)
	if err != nil {
		return err
	}
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

	tmpConfig.SetDefaults()
	s.Config = tmpConfig
	return nil
}

func (s *Server) HandleDiscovery(req *api.DiscoveryRequest) (*api.DiscoveryResponse, error) {
	var resp api.DiscoveryResponse
	return &resp, nil
}

func (s *Server) HandleConfig(req *api.ConfigRequest) (*api.ConfigResponse, error) {
	// Copy the public parts of our configuration.
	var resp api.ConfigResponse
	resp.Config = api.ConfigPublic{
		Name:     s.Config.Name,
		Idp:      s.Config.Idp,
		Roles:    s.Config.Roles,
		Workflow: s.Config.Workflow,
	}
	return &resp, nil
}

func (s *Server) HandleDirectSamlAuth(eq *api.DirectSamlAuthRequest) (*api.DirectAuthResponse, error) {
	return nil, errors.New("Not implemented")
}

func (s *Server) HandleDirectOidcAuth(req *api.DirectOidcAuthRequest) (*api.DirectAuthResponse, error) {
	return nil, errors.New("Not implemented")
}

func (s *Server) HandleWorkflowStart(req *api.WorkflowStartRequest) (*api.WorkflowStartResponse, error) {
	// NOTE: this will be a JWT in future to allow stateless operation.
	uu := uuid.New()
	return &api.WorkflowStartResponse{
		Nonce: uu.String(),
	}, nil
}

func (s *Server) HandleWorkflowAuth(req *api.WorkflowAuthRequest) (*api.WorkflowAuthResponse, error) {
	role := s.Config.FindRoleByName(req.Role)
	if role == nil {
		return nil, errors.Errorf("requested role not found: %s", req.Role)
	}
	credIssuer, err := creds.NewFromConfig(role, &s.Config)

	if err != nil {
		return nil, errors.Wrap(err, "during issuer configuration")
	}
	// TODO: check SAML assertions here

	userInfo := api.AuthInfo{
		Environment: s.Config.Name,
		Role: req.Role,
		Username: req.Username,
		ValidFor: role.ValidForSeconds,
	}
	issuedCreds, err := credIssuer.IssueFor(&userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "during issuance")
	}
	return &api.WorkflowAuthResponse{
		Credentials: issuedCreds,
	}, nil
}
