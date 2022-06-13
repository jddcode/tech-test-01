package weatherFetcher

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jddcode/tech-test-ennismore/internal/mocks"
	"github.com/jddcode/tech-test-ennismore/internal/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Tests")
}

var _ = Describe("Weather forecast handler", func() {
	var (
		mockController *gomock.Controller
		mockHttpClient *mocks.MockClient
		mockFetcher    weatherFetcher
	)

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
		mockHttpClient = mocks.NewMockClient(mockController)
		mockFetcher = weatherFetcher{
			web: mockHttpClient,
		}
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Context("Fetching a weather forecast for a specific location", func() {
		When("the initial lat/long based GET request fails", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return("", errors.New("some http error"))
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(fmt.Errorf(ErrorGetRequest, "some http error")))
			})
		})

		When("the initial lat/long based GET request response cannot be unmarshalled", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return("---", nil)
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(fmt.Errorf(ErrorUnmarshalLookup, "invalid character '-' in numeric literal")))
			})
		})

		When("the initial lat/long based GET request gives a blank or unpopulated forecast URL", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":""}}`, nil)
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(errors.New(ErrorNoForecastResource)))
			})
		})

		When("the forecast URL has an error during the GET request", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":"http://example.org"}}`, nil)
				mockHttpClient.EXPECT().Get("http://example.org").Return("", errors.New("some http error"))
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(fmt.Errorf(ErrorGetForecast, "some http error")))
			})
		})

		When("the data received from the forecast lookup cannot be unmarshalled", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":"http://example.org"}}`, nil)
				mockHttpClient.EXPECT().Get("http://example.org").Return("---", nil)
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(fmt.Errorf(ErrorUnmarshalForecast, "invalid character '-' in numeric literal")))
			})
		})

		When("the data received from the forecast lookup contains an invalid start time", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":"http://example.org"}}`, nil)
				mockHttpClient.EXPECT().Get("http://example.org").Return(`{"properties":{"periods":[{"startTime":"invalid"}]}}`, nil)
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(fmt.Errorf(ErrorUnusualStartTime, "invalid time string")))
			})
		})

		When("the data received from the forecast lookup contains an invalid end time", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":"http://example.org"}}`, nil)
				mockHttpClient.EXPECT().Get("http://example.org").Return(`{"properties":{"periods":[{"startTime":"2022-01-01T13:00:00", "endTime":"invalid"}]}}`, nil)
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(fmt.Errorf(ErrorUnusualEndTime, "invalid time string")))
			})
		})

		When("the data received from the forecast lookup contains an invalid wind speed", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":"http://example.org"}}`, nil)
				mockHttpClient.EXPECT().Get("http://example.org").Return(`{"properties":{"periods":[{"startTime":"2022-01-01T13:00:00", "endTime":"2022-01-01T18:00:00", "windSpeed": "invalid"}]}}`, nil)
				_, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).To(Equal(fmt.Errorf(ErrorUnusualWindSpeed, "Unexpected format for wind speed string")))
			})
		})

		When("the data received is valid and there is an upper and lower wind speed", func() {
			It("should return a slice of weather forecasts", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":"http://example.org"}}`, nil)
				mockHttpClient.EXPECT().Get("http://example.org").Return(
					`{"properties":{"periods":[{"startTime":"2022-01-01T13:00:00", "endTime":"2022-01-01T18:00:00", "windSpeed": "4 to 8 mph", "shortForecast": "it will be sunny"}]}}`, nil)
				predictions, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).ToNot(HaveOccurred())
				Expect(predictions[0].Wind.MinSpeed).To(Equal(4))
				Expect(predictions[0].Wind.MaxSpeed).To(Equal(8))
				Expect(predictions[0].GetForecast()).To(Equal("it will be sunny"))
			})
		})

		When("the data received is valid and there is only one wind speed", func() {
			It("should return a slice of weather forecasts", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`{"properties":{"forecast":"http://example.org"}}`, nil)
				mockHttpClient.EXPECT().Get("http://example.org").Return(
					`{"properties":{"periods":[{"startTime":"2022-01-01T13:00:00", "endTime":"2022-01-01T18:00:00", "windSpeed": "5 mph", "shortForecast": "it will be sunny"}]}}`, nil)
				predictions, err := mockFetcher.Fetch(structs.CoOrdinates{})

				Expect(err).ToNot(HaveOccurred())
				Expect(predictions[0].Wind.MinSpeed).To(Equal(5))
				Expect(predictions[0].Wind.MaxSpeed).To(Equal(5))
				Expect(predictions[0].GetForecast()).To(Equal("it will be sunny"))
			})
		})
	})
})
