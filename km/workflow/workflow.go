package workflow

import "context"

type Workflow struct {
	BaseUrl string
	PolicyEncrypter string // TODO: key, signingmethod etc
}

type StartRequest struct {
	Requester struct {
		Name string `json:"name"`
		Username string `json:"username"`
		Email string `json:"email"`
	} `json:"requester"`
	Source struct {
		// Source Key? For gitlab?
		Description string `json:"description"`
		DetailsURI string `json:"details_uri"`
	} `json:"source"`
	Target struct {
		Name string `json:"name"`
		// Target discovery URI?
	} `json:"target"`
	// To be encrypted with workflow public key
	WorkflowPolicy string `json:"workflow_policy"`
}

type StartResponse struct {
	// Workflow id
	// Workflow nonce
}

type GetDetailsRequest struct {
	// Workflow id
}

type GetDetailsResponse struct {
	// State: pending / completed
}

type GetAssertionsRequest struct {
	// Workflow id
	// Workflow nonce
}

type GetAssertionsResponse struct {
	// Bag of SAML assertions. Could be wrapped(?)
}

func (w *Workflow) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	// HTTP
	return nil, nil
}

func (w *Workflow) GetDetails(ctx context.Context, req *GetDetailsRequest) (*GetDetailsResponse, error) {
	// HTTP
	return nil, nil
}

func (w *Workflow) GetAssertions(ctx context.Context, req *GetAssertionsRequest) (*GetAssertionsResponse, error) {
	return nil, nil
}
