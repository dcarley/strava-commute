package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

const ChallengeKey = "hub.challenge"

type CallbackResponse struct {
	Challenge string `json:"hub.challenge"`
}

func CallbackHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	challenge := request.QueryStringParameters[ChallengeKey]
	if challenge == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("missing query param: %s\n", ChallengeKey),
		}
	}

	resp, err := json.Marshal(CallbackResponse{
		Challenge: challenge,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       fmt.Sprintf("%s\n", err),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("%s\n", resp),
	}
}
