package main_test

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/dcarley/strava-commute"
)

var _ = Describe("Main", func() {
	const exampleBody = `{
	"subscription_id": "1",
	"owner_id": 13408,
	"object_id": 12312312312,
	"object_type": "activity",
	"aspect_type": "create",
	"events_time": 1297286541
}`

	DescribeTable("Handler",
		func(request events.APIGatewayProxyRequest, expected events.APIGatewayProxyResponse) {
			resp, err := Handler(request)
			Expect(err).ToNot(HaveOccurred(), "error should be nil for API Gateway to send response")
			Expect(resp).To(Equal(expected))
		},
		Entry("invalid request method",
			events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusMethodNotAllowed,
				Body:       ``,
			},
		),
		Entry("empty request body",
			events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       ``,
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `unexpected end of JSON input`,
			},
		),
		Entry("invalid request body",
			events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       `{"key": valyou}`,
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `invalid character 'v' looking for beginning of value`,
			},
		),
		Entry("valid request body",
			events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       exampleBody,
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       ``,
			},
		),
	)
})
