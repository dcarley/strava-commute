package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
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

func UpdateHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	token := os.Getenv(TokenEnvVar)
	if token == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       fmt.Sprintf("environment variable not set: %s\n", TokenEnvVar),
		}
	}

	config, err := LoadConfig(ConfigFile)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       fmt.Sprintf("%s\n", err),
		}
	}

	var webhook WebhookRequest
	err = json.Unmarshal([]byte(request.Body), &webhook)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("%s\n", err),
		}
	}

	client := NewClient(token)
	service := strava.NewActivitiesService(client)
	activity, err := service.Get(webhook.ObjectID).Do()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       fmt.Sprintf("%s\n", err),
		}
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
			Body:       fmt.Sprintf("no location matches for: %d\n", activity.Id),
		}
	}

	name := fmt.Sprintf("Commute%s%s", startName, endName)
	update := service.Update(activity.Id).Name(name).Commute(true)
	if id := config.GearID; id != "" {
		update.Gear(id)
	}
	_, err = update.Do()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       fmt.Sprintf("%s\n", err),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("renamed %d to: %s\n", activity.Id, name),
	}
}
