package main_test

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/dcarley/strava-commute"
)

var _ = Describe("CallbackHandler", func() {
	DescribeTable("callbacks",
		func(request events.APIGatewayProxyRequest, expected events.APIGatewayProxyResponse) {
			Expect(CallbackHandler(request)).To(Equal(expected))
		},
		Entry("success",
			events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
				QueryStringParameters: map[string]string{
					"hub.mode":      "subscribe",
					"hub.challenge": "mychallenge",
				},
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       "{\"hub.challenge\":\"mychallenge\"}\n",
			},
		),
		Entry("missing challenge",
			events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
				QueryStringParameters: map[string]string{
					"hub.mode": "subscribe",
				},
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "missing query param: hub.challenge\n",
			},
		),
	)
})
