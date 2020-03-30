package api

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	"log"
)

type Client struct {
	//    * Function name - my-function (name-only), my-function:v1 (with alias).
	//    * Function ARN - arn:aws:lambda:us-west-2:123456789012:function:my-function.
	//    * Partial ARN - 123456789012:function:my-function.
	FunctionName string
	lambdaClient *lambda.Lambda
}

func NewClient(target string) *Client {
	c := new(Client)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	c.FunctionName = target
	c.lambdaClient = lambda.New(sess) // TODO: region? Or can that come from env?
	return c
}

func (c *Client) GetConfig(req *ConfigRequest) (*ConfigResponse, error) {
	resp := new(ConfigResponse)
	err := c.rpc(&Request{ Type: "config", Payload: req}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DirectSamlAuth(req *DirectSamlAuthRequest) (*DirectAuthResponse, error) {
	return nil, nil
}

func (c *Client) WorkflowStart(req *WorkflowStartRequest) (*WorkflowStartResponse, error) {
	resp := new(WorkflowStartResponse)
	err := c.rpc(&Request{ Type: "workflow_start", Payload: req}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) WorkflowAuth(req *WorkflowAuthRequest) (*WorkflowAuthResponse, error) {
	resp := new(WorkflowAuthResponse)
	err := c.rpc(&Request{ Type: "workflow_auth", Payload: req}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) rpc(req interface{}, resp interface{}) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "rpc marshal")
	}
	log.Println("km lambda request:" + string(payload))
	result, err := c.lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(c.FunctionName),
		Payload: payload,
	})
	if err != nil {
		return errors.Wrap(err, "rpc lambda invoke")
	}
	// TODO: think about other stuff in invoke response
	log.Println("km lambda response:" + string(result.Payload))
	err = json.Unmarshal(result.Payload, resp)
	if err != nil {
		return errors.Wrap(err, "rpc unmarshal")
	}
	return nil
}