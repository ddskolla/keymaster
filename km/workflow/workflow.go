package workflow

import "context"

type Workflow struct {

}

type StartRequest struct {

}

type StartResponse struct {

}

type GetDetailsRequest struct {

}

type GetDetailsResponse struct {
	
}

type GetAssertionsRequest struct {

}

type GetAssertionsResponse struct {

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
