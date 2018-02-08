package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	strava "github.com/strava/go.strava"
)

const (
	TokenEnvVar = "STRAVA_API_TOKEN"
	ConfigFile  = "config.json"
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
	token := os.Getenv(TokenEnvVar)
	if token == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       fmt.Sprintf("%s environment variable not set", TokenEnvVar),
		}, nil
	}

	config, err := LoadConfig(ConfigFile)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       err.Error(),
		}, nil
	}

	if request.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}

	var webhook WebhookRequest
	err = json.Unmarshal([]byte(request.Body), &webhook)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	client := NewClient(token)
	service := strava.NewActivitiesService(client)
	activity, err := service.Get(webhook.ObjectID).Do()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       err.Error(),
		}, nil
	}

	var startName, endName string
	if name := config.GetLocation(activity.StartLocation); name != "" {
		startName = fmt.Sprintf(" from %s", name)
	}
	if name := config.GetLocation(activity.EndLocation); name != "" {
		endName = fmt.Sprintf(" to %s", name)
	}
	if startName == "" && endName == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       fmt.Sprintf("no need to rename %d", activity.Id),
		}, nil
	}

	name := fmt.Sprintf("Commute%s%s", startName, endName)
	_, err = service.Update(activity.Id).Name(name).Commute(true).Do()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("renamed %d to: %s", activity.Id, name),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
