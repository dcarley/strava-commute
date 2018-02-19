package main_test

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dcarley/strava-commute"
)

var _ = Describe("Handler", func() {
	Describe("invalid request method", func() {
		It("should return 405", func() {
			resp, err := Handler(events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
			})
			Expect(resp).To(Equal(events.APIGatewayProxyResponse{
				StatusCode: http.StatusMethodNotAllowed,
				Body:       ``,
			}))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
