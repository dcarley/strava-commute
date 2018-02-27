package main_test

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	. "github.com/dcarley/strava-commute"
)

var _ = Describe("Handler", func() {
	var logBuffer *gbytes.Buffer

	BeforeEach(func() {
		logBuffer = gbytes.NewBuffer()
		Logger = log.New(logBuffer, "", 0)
	})

	DescribeTable("request methods",
		func(req events.APIGatewayProxyRequest, expected events.APIGatewayProxyResponse, logs string) {
			resp, err := Handler(req)
			Expect(resp).To(Equal(expected))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(logBuffer.Contents())).To(Equal(logs))
		},
		Entry("should log error response",
			events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusMethodNotAllowed,
				Body:       "unsupported method: PUT\n",
			},
			"Error: 405 unsupported method: PUT\n",
		),
		Entry("not log successful response",
			events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
				QueryStringParameters: map[string]string{
					"hub.challenge": "foo",
				},
			},
			events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       "{\"hub.challenge\":\"foo\"}\n",
			},
			"",
		),
	)
})
