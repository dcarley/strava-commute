package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// error is always nil, otherwise API gateway won't deliver the HTTP
// response with the correct status code and body.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodPost:
		return UpdateHandler(request), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
