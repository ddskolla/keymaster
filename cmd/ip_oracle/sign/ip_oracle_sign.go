package main

import (
	"encoding/json"
    "net/http"

    "github.com/bsycorp/keymaster/km/ip_oracle"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

type RequestPayLoad struct {
    KMSKeyId string
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var reqPayLoad RequestPayLoad
    json.Unmarshal([]byte(req.Body), &reqPayLoad)
    
    sm := ip_oracle.NewSigningMethodKMS(reqPayLoad.KMSKeyId)
    signedString, err := ip_oracle.MakeIPJWT(req.RequestContext.Identity.SourceIP, sm)

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       signedString,
    }, err
}

func main() {
    lambda.Start(HandleRequest)
}
