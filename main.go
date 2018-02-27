package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var Logger = log.New(os.Stderr, "", log.LstdFlags)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var resp events.APIGatewayProxyResponse
	switch request.HTTPMethod {
	case http.MethodGet:
		resp = CallbackHandler(request)
	case http.MethodPost:
		resp = UpdateHandler(request)
	default:
		resp = events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       fmt.Sprintf("unsupported method: %s\n", request.HTTPMethod),
		}
	}

	if resp.StatusCode >= http.StatusBadRequest {
		Logger.Printf("Error: %d %s", resp.StatusCode, resp.Body)
	}

	// error is always nil, otherwise API gateway won't deliver the HTTP
	// response with the correct status code and body.
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
