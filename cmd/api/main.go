package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/pkg/errors"
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

func Handler(ctx context.Context, req api.Request) (interface{}, error) {
	switch r := req.Payload.(type) {
	case *api.ConfigRequest:
		return HandleConfig(ctx, r)
	case *api.DirectSamlAuthRequest:
		return HandleDirectSamlAuth(ctx, r)
	case *api.DirectOidcAuthRequest:
		return HandleDirectOidcAuth(ctx, r)
	case *api.WorkflowStartRequest:
		return HandleWorkflowStart(ctx, r)
	case *api.WorkflowAuthRequest:
		return HandleWorkflowAuth(ctx, r)
	default:
		return nil, errors.New("unexpected request")
	}
}

func main() {
	lambda.Start(Handler)
}
