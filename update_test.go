package main_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

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

var _ = Describe("UpdateHandler", func() {
	const (
		activityID    = 12312312312
		eventTemplate = `{
				"subscription_id": "1",
				"owner_id": 13408,
				"object_id": %d,
				"object_type": "activity",
				"aspect_type": "create",
				"events_time": 1297286541
			}`
	)

	var (
		pwd, tempDir, configFile string
		server                   *ghttp.Server
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

		tempDir, err := ioutil.TempDir("", "strava")
		if err != nil {
			log.Fatal(err)
		}

		basicConfig := []byte(`{
			"locations": {
				"null": {}
			}
		}`)
		configFile = path.Join(tempDir, ConfigFile)
		Expect(ioutil.WriteFile(configFile, basicConfig, 0644)).To(Succeed())

		pwd, err = os.Getwd()
		Expect(err).ToNot(HaveOccurred())
		Expect(os.Chdir(tempDir)).To(Succeed())
	})

	AfterEach(func() {
		server.Close()
		Expect(os.Unsetenv(TokenEnvVar)).To(Succeed())

		Expect(os.Chdir(pwd)).To(Succeed())
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	Describe("token environment variable not set", func() {
		BeforeEach(func() {
			Expect(os.Unsetenv(TokenEnvVar)).To(Succeed())
		})

		It("should respond with 503", func() {
			Expect(UpdateHandler(events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       `{}`,
			})).To(Equal(events.APIGatewayProxyResponse{
				StatusCode: http.StatusServiceUnavailable,
				Body:       "environment variable not set: STRAVA_API_TOKEN\n",
			}))
		})
	})

	Describe("config", func() {
		Describe("config file doesn't exist", func() {
			BeforeEach(func() {
				Expect(os.Remove(configFile)).To(Succeed())
			})

			It("should respond with 503", func() {
				Expect(UpdateHandler(events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Body:       `{}`,
				})).To(Equal(events.APIGatewayProxyResponse{
					StatusCode: http.StatusServiceUnavailable,
					Body:       "open config.json: no such file or directory\n",
				}))
			})
		})

		Describe("config file contains no locations", func() {
			BeforeEach(func() {
				config := []byte(`{
					"locations": {}
				}`)
				Expect(ioutil.WriteFile(configFile, config, 0644)).To(Succeed())
			})

			It("should respond with 503", func() {
				Expect(UpdateHandler(events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Body:       `{}`,
				})).To(Equal(events.APIGatewayProxyResponse{
					StatusCode: http.StatusServiceUnavailable,
					Body:       "config contains no locations\n",
				}))
			})
		})
	})

	Describe("does not fetch activity from Strava API", func() {
		AfterEach(func() {
			Expect(server.ReceivedRequests()).To(HaveLen(0))
		})

		DescribeTable("invalid requests",
			func(request events.APIGatewayProxyRequest, expected events.APIGatewayProxyResponse) {
				Expect(UpdateHandler(request)).To(Equal(expected))
			},
			Entry("empty request body",
				events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Body:       ``,
				},
				events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "unexpected end of JSON input\n",
				},
			),
			Entry("invalid request body",
				events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Body:       `{"key": valyou}`,
				},
				events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "invalid character 'v' looking for beginning of value\n",
				},
			),
		)
	})

	Describe("location handling", func() {
		BeforeEach(func() {
			config := []byte(`{
				"locations": {
					"London": {
						"min": [51.286758, -0.510375],
						"max": [51.691875, 0.334015]
					},
					"Sheffield": {
						"min": [53.304512, -1.801472],
						"max": [53.503128, -1.324669]
					}
				}
			}`)
			Expect(ioutil.WriteFile(configFile, config, 0644)).To(Succeed())
		})

		Describe("no location matched", func() {
			BeforeEach(func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/v3/activities/%d", activityID)),
					ghttp.RespondWithJSONEncoded(http.StatusOK, strava.ActivityDetailed{
						ActivitySummary: strava.ActivitySummary{
							Id:            activityID,
							Name:          "Morning Ride",
							StartLocation: strava.Location{25, 25},
							EndLocation:   strava.Location{25, 25},
						},
					}),
				))
			})

			It("should not rename activity", func() {
				Expect(UpdateHandler(events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Body:       fmt.Sprintf(eventTemplate, activityID),
				})).To(Equal(events.APIGatewayProxyResponse{
					StatusCode: http.StatusOK,
					Body:       fmt.Sprintf("no location matches for: %d\n", activityID),
				}))
			})
		})

		DescribeTable("start and end locations",
			func(start, end strava.Location, name string) {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("/api/v3/activities/%d", activityID)),
						ghttp.RespondWithJSONEncoded(http.StatusOK, strava.ActivityDetailed{
							ActivitySummary: strava.ActivitySummary{
								Id:            activityID,
								Name:          "Morning Ride",
								StartLocation: start,
								EndLocation:   end,
							},
						}),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT",
							fmt.Sprintf("/api/v3/activities/%d", activityID),
							fmt.Sprintf("name=%s&commute=true", name),
						),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
				)

				Expect(UpdateHandler(events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Body:       fmt.Sprintf(eventTemplate, activityID),
				})).To(Equal(events.APIGatewayProxyResponse{
					StatusCode: http.StatusOK,
					Body:       fmt.Sprintf("renamed %d to: %s\n", activityID, name),
				}))
				Expect(server.ReceivedRequests()).To(HaveLen(2))
			},
			Entry("start matches",
				strava.Location{51.5, 0},
				strava.Location{25, 25},
				"Commute from London",
			),
			Entry("end matches",
				strava.Location{25, 25},
				strava.Location{53.5, -1.5},
				"Commute to Sheffield",
			),
			Entry("start and end matches",
				strava.Location{51.5, 0},
				strava.Location{53.5, -1.5},
				"Commute from London to Sheffield",
			),
		)
	})

	Describe("gear ID in config", func() {
		const (
			name           = "Commute from simple"
			gearID         = "12345"
			configTemplate = `{
				"gear_id": "%s",
				"locations": {
					"simple": {
						"min": [0, 0],
						"max": [1, 1]
					}
				}
			}`
		)

		BeforeEach(func() {
			config := []byte(fmt.Sprintf(configTemplate, gearID))
			Expect(ioutil.WriteFile(configFile, config, 0644)).To(Succeed())

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/v3/activities/%d", activityID)),
					ghttp.RespondWithJSONEncoded(http.StatusOK, strava.ActivityDetailed{
						ActivitySummary: strava.ActivitySummary{
							Id:            activityID,
							Name:          "Morning Ride",
							StartLocation: strava.Location{0.5, 0.5},
							EndLocation:   strava.Location{25, 25},
						},
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT",
						fmt.Sprintf("/api/v3/activities/%d", activityID),
						fmt.Sprintf("name=%s&commute=true&gear_id=%s", name, gearID),
					),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
			)
		})

		AfterEach(func() {
			Expect(server.ReceivedRequests()).To(HaveLen(2))
		})

		It("should set gear ID for matching activity", func() {
			Expect(UpdateHandler(events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       fmt.Sprintf(eventTemplate, activityID),
			})).To(Equal(events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       fmt.Sprintf("renamed %d to: %s\n", activityID, name),
			}))
		})
	})
})
