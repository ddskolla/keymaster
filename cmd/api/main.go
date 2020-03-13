package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"

	api "github.com/bsycorp/keymaster/keymaster/api"
)

func HandlePing(ctx context.Context, req api.PingRequestMessage) (interface{}, error) {
	return api.PingResponseMessage{Message:"OK"}, nil
}

func Handler(ctx context.Context, reqBytes []byte) (interface{}, error) {
	var req api.Request
	if err := json.Unmarshal(reqBytes, &req); err != nil {
		return nil, err
	}
	switch msg := req.Message.(type) {
	case api.PingRequestMessage:
		return HandlePing(ctx, msg)
	}
	return nil, errors.New("unhandled operation")
}

func main() {
	lambda.Start(Handler)
}
