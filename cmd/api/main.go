package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bsycorp/keymaster/keymaster/api"
	"github.com/pkg/errors"
	"os"
)

func HandleConfig(ctx context.Context, req *api.ConfigRequest) (*api.ConfigResponse, error) {
	return nil, errors.New("Not implemented")
}

func HandleDirectSamlAuth(ctx context.Context, req *api.DirectSamlAuthRequest) (*api.DirectAuthResponse, error) {
	return nil, errors.New("Not implemented")
}

func HandleDirectOidcAuth(ctx context.Context, req *api.DirectOidcAuthRequest) (*api.DirectAuthResponse, error) {
	return nil, errors.New("Not implemented")
}

func HandleWorkflowStart(ctx context.Context, req *api.WorkflowStartRequest) (*api.WorkflowStartResponse, error) {
	return nil, errors.New("Not implemented")
}

func HandleWorkflowAuth(ctx context.Context, req *api.WorkflowAuthRequest) (*api.WorkflowAuthResponse, error) {
	return nil, errors.New("Not implemented")
}

func main() {
	switch os.Getenv("_HANDLER") {
	case "config":
		lambda.Start(HandleConfig)
	case "direct_saml_auth":
		lambda.Start(HandleDirectSamlAuth)
	case "direct_oidc_auth":
		lambda.Start(HandleDirectOidcAuth)
	case "workflow_start":
		lambda.Start(HandleWorkflowStart)
	case "workflow_auth":
		lambda.Start(HandleWorkflowAuth)
	}
	// TODO default case?
}
