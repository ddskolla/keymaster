package api

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type Client struct {
	//    * Function name - my-function (name-only), my-function:v1 (with alias).
	//    * Function ARN - arn:aws:lambda:us-west-2:123456789012:function:my-function.
	//    * Partial ARN - 123456789012:function:my-function.
	FunctionName string
	lambdaClient *lambda.Lambda
}

func NewClient() *Client {
	c := new(Client)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	c.lambdaClient = lambda.New(sess) // TODO: region? Or can that come from env?
	return c
}

func (c *Client) GetConfig(req *ConfigRequest) (*ConfigResponse, error) {
	return nil, nil
}

func (c *Client) DirectSamlAuth(req *DirectSamlAuthRequest) (*DirectAuthResponse, error) {
	return nil, nil
}

func (c *Client) WorkflowStart(req *WorkflowStartRequest) (*WorkflowStartResponse, error) {
	return nil, nil
}

func (c *Client) WorkflowAuth(req *WorkflowAuthRequest) (*WorkflowAuthResponse, error) {
	return nil, nil
}

func (c *Client) rpc(req interface{}) ([]byte, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	result, err := c.lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String("MyGetItemsFunction"),
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}
	// TODO: looks like there is more stuff in the result to look at
	return result.Payload, nil
}