package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL *url.URL
	HttpClient *http.Client
	PolicyEncrypter string // TODO: key, signingmethod etc
}

func NewClient(baseUrl string) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing workflow client URL")
	}
	return &Client{
		BaseURL:         u,
		HttpClient:      http.DefaultClient,
		PolicyEncrypter: "", // TODO
	}, nil
}

type Requester struct {
	Name string `json:"name"`
	Username string `json:"username"`
	Email string `json:"email"`
}

type Source struct {
	// Source Key? For gitlab?
	Description string `json:"description"`
	DetailsURI string `json:"details_uri"`
}

type Target struct {
	EnvironmentName string `json:"environment_name"`
	EnvironmentDiscoveryURI string `json:"environment_discovery_uri"`
}

type Policy struct {
	Name                string         `json:"name"`
	IdpName             string         `json:"idp_name"`
	RequesterCanApprove bool           `json:"requester_can_approve"`
	IdentifyRoles       map[string]int `json:"identify_roles"`
	ApproverRoles       map[string]int `json:"approver_roles"`
}

type StartRequest struct {
	Requester Requester `json:"requester"`
	Source Source `json:"source"`
	Target Target`json:"target"`
	// To be encrypted with workflow public key
	// TODO: change type here, separation of concerns...
	Policy Policy `json:"policy"`
}

type StartResponse struct {
	WorkflowId string `json:"workflow_id"`
	WorkflowUrl string `json:"workflow_url"`
	Nonce string `json:"nonce"`
}

type GetDetailsRequest struct {
	// Workflow id
}

type GetDetailsResponse struct {
	// State: pending / completed
}

type GetAssertionsRequest struct {
	WorkflowId string `json:"workflow_id"`
	Nonce string `json:"nonce"`
}

type GetAssertionsResponse struct {
	// Bag of SAML assertions. Could be wrapped(?)
	//Workflow Workflow `json:"workflow"` //???
	Status string `json:"status"` //???
	Assertions []string `json:"assertions"` // Resulting IDP assertions
}

func (c *Client) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	httpReq, err := c.newRequest(ctx, "POST", "/1/workflow/create", req)
	if err != nil {
		return nil, err
	}
	var resp StartResponse
	_, err = c.do(httpReq, &resp)
	return &resp, err
}

func (c *Client) GetDetails(ctx context.Context, req *GetDetailsRequest) (*GetDetailsResponse, error) {
	httpReq, err := c.newRequest(ctx, "POST", "/1/workflow/getDetails", req)
	if err != nil {
		return nil, err
	}
	var resp GetDetailsResponse
	_, err = c.do(httpReq, &resp)
	return &resp, err
}

func (c *Client) GetAssertions(ctx context.Context, req *GetAssertionsRequest) (*GetAssertionsResponse, error) {
	httpReq, err := c.newRequest(ctx, "POST", "/1/workflow/getAssertions", req)
	if err != nil {
		return nil, err
	}
	var resp GetAssertionsResponse
	_, err = c.do(httpReq, &resp)
	return &resp, err
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading workflow response body")
	}
	log.Printf("RESPONSE FROM WORKFLOW: %s", string(body))
	err = json.Unmarshal(body, v)
	return resp, err
}
