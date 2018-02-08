package main_test

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	strava "github.com/strava/go.strava"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/dcarley/strava-commute"
)

type MockTransport struct {
	Host      string
	Transport *http.Transport
}

func (m MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = m.Host

	return m.Transport.RoundTrip(req)
}

var _ = Describe("Main", func() {
	const (
		errorDescription = "error should be nil for API Gateway to send response"
	)

	Describe("token environment variable not set", func() {
		BeforeEach(func() {
			Expect(os.Unsetenv(TokenEnvVar)).To(Succeed())
		})

		It("should respond with 503", func() {
			resp, err := Handler(events.APIGatewayProxyRequest{})
			Expect(err).ToNot(HaveOccurred(), errorDescription)
			Expect(resp).To(Equal(events.APIGatewayProxyResponse{
				StatusCode: http.StatusServiceUnavailable,
				Body:       `STRAVA_API_TOKEN environment variable not set`,
			}))
		})
	})

	Describe("token envirnment variable set", func() {
		var (
			server *ghttp.Server
		)

		BeforeEach(func() {
			Expect(os.Setenv(TokenEnvVar, "mytoken")).To(Succeed())

			server = ghttp.NewServer()
			serverURL, err := url.Parse(server.URL())
			Expect(err).ToNot(HaveOccurred())
			Transport = MockTransport{
				Host:      serverURL.Host,
				Transport: &http.Transport{},
			}
		})

		AfterEach(func() {
			server.Close()
			Expect(os.Unsetenv(TokenEnvVar)).To(Succeed())
		})

		Describe("does not fetch activity from Strava API", func() {
			AfterEach(func() {
				Expect(server.ReceivedRequests()).To(HaveLen(0))
			})

			DescribeTable("invalid requests",
				func(request events.APIGatewayProxyRequest, expected events.APIGatewayProxyResponse) {
					resp, err := Handler(request)
					Expect(err).ToNot(HaveOccurred(), errorDescription)
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
			)
		})

		Describe("does fetch activity from Strava API", func() {
			const (
				objectID            = 12312312312
				requestBodyTemplate = `{
	"subscription_id": "1",
	"owner_id": 13408,
	"object_id": %d,
	"object_type": "activity",
	"aspect_type": "create",
	"events_time": 1297286541
}`
			)

			DescribeTable("valid requests",
				func(response strava.ActivitySummary, expected events.APIGatewayProxyResponse) {
					server.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("/api/v3/activities/%d", objectID)),
						ghttp.RespondWithJSONEncoded(http.StatusOK, response),
					))

					resp, err := Handler(events.APIGatewayProxyRequest{
						HTTPMethod: "POST",
						Body:       fmt.Sprintf(requestBodyTemplate, objectID),
					})
					Expect(err).ToNot(HaveOccurred(), errorDescription)
					Expect(resp).To(Equal(expected))
				},
				Entry("basic test name",
					strava.ActivitySummary{
						Name: "test",
					},
					events.APIGatewayProxyResponse{
						StatusCode: http.StatusOK,
						Body:       `test`,
					},
				),
			)
		})
	})
})
