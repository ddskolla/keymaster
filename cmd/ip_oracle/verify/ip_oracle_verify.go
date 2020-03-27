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
    SignedString string
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var reqPayLoad RequestPayLoad
    json.Unmarshal([]byte(req.Body), &reqPayLoad)

    sm := ip_oracle.NewSigningMethodKMS(reqPayLoad.KMSKeyId)
    verifiedSourceIp, err := ip_oracle.VerifyIPJWT(reqPayLoad.SignedString, sm)

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       verifiedSourceIp,
    }, err
}

func main() {
    lambda.Start(HandleRequest)
}
