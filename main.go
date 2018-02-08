package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// http://strava.github.io/api/partner/v3/events/#updates
type WebhookRequest struct {
	SubscriptionID string                  `json:"subscription_id"`
	OwnerID        int64                   `json:"owner_id"`
	ObjectID       int64                   `json:"object_id"`
	ObjectType     string                  `json:"object_type"`
	AspectType     string                  `json:"aspect_type"`
	EventTime      events.SecondsEpochTime `json:"event_time"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}

	var webhook WebhookRequest
	err := json.Unmarshal([]byte(request.Body), &webhook)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
