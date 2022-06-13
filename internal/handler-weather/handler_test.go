package handlerWeather

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	handlerStructs "github.com/jddcode/tech-test-ennismore/internal/handler-weather/structs"
	"github.com/jddcode/tech-test-ennismore/internal/mocks"
	"github.com/jddcode/tech-test-ennismore/internal/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Tests")
}

var _ = Describe("Weather forecast handler", func() {
	var (
		mockController     *gomock.Controller
		mockCoordinates    *mocks.MockFinder
		mockWeatherFetcher *mocks.MockWeatherFetcher
		mockCache          *mocks.MockCache
		mockHandler        handler
	)

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
		mockCoordinates = mocks.NewMockFinder(mockController)
		mockWeatherFetcher = mocks.NewMockWeatherFetcher(mockController)
		mockCache = mocks.NewMockCache(mockController)
		mockHandler = handler{
			coOrdinates: mockCoordinates,
			weather:     mockWeatherFetcher,
			cache:       mockCache,
		}
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Context("Requesting an update on the weather", func() {
		When("a request is received with no cities", func() {
			It("should return an error", func() {
				mockReq, _ := http.NewRequest(http.MethodGet, "/weather", nil)
				resp := httptest.NewRecorder()
				mockHandler.Handle(resp, mockReq)

				result := resp.Result()
				defer result.Body.Close()
				data, err := ioutil.ReadAll(result.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(data)).To(Equal(ErrorNoCities))
			})
		})

		When("a request is received with a city we cannot get co-ordinates for", func() {
			It("should return an error", func() {
				mockCache.EXPECT().Get("testcity").Return(nil, errors.New("cache miss"))
				mockReq, _ := http.NewRequest(http.MethodGet, "/weather?city=testcity", nil)
				resp := httptest.NewRecorder()

				mockCoordinates.EXPECT().Find("testcity", "usa").Return(structs.CoOrdinates{}, errors.New("could not find co-ordinates"))
				mockHandler.Handle(resp, mockReq)

				result := resp.Result()
				defer result.Body.Close()
				data, err := ioutil.ReadAll(result.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(data)).To(Equal(fmt.Sprintf(ErrorNoCoordinates, "testcity")))
			})
		})

		When("a request is received we cannot get a forecast for", func() {
			It("should return an error", func() {
				mockCache.EXPECT().Get("testcity").Return(nil, errors.New("cache miss"))
				mockReq, _ := http.NewRequest(http.MethodGet, "/weather?city=testcity", nil)
				resp := httptest.NewRecorder()

				mockCoordinates.EXPECT().Find("testcity", "usa").Return(structs.CoOrdinates{}, nil)
				mockWeatherFetcher.EXPECT().Fetch(structs.CoOrdinates{}).Return([]structs.Weather{structs.Weather{}}, errors.New("could not fetch forecast"))
				mockHandler.Handle(resp, mockReq)

				result := resp.Result()
				defer result.Body.Close()
				data, err := ioutil.ReadAll(result.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(data)).To(Equal(fmt.Sprintf(ErrorNoForecast, "testcity")))
			})
		})

		When("everything is working", func() {
			It("should return a json weather forecast", func() {
				mockCache.EXPECT().Get("testcity").Return(nil, errors.New("cache miss"))
				mockReq, _ := http.NewRequest(http.MethodGet, "/weather?city=testcity", nil)
				resp := httptest.NewRecorder()

				mockCoordinates.EXPECT().Find("testcity", "usa").Return(structs.CoOrdinates{}, nil)

				setTime, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 12:00:00")
				weatherResult := structs.Weather{
					Start: setTime,
					End:   setTime,
				}
				weatherResult.Forecast.Long = "long dry spells"
				mockWeatherFetcher.EXPECT().Fetch(structs.CoOrdinates{}).Return([]structs.Weather{weatherResult}, nil)

				mockCache.EXPECT().Store("testcity", gomock.Any())

				mockHandler.Handle(resp, mockReq)

				result := resp.Result()
				defer result.Body.Close()
				data, err := ioutil.ReadAll(result.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(data)).To(Equal(`{"forecast":[{"name":"testcity","detail":[{"starttime":"2020-01-01T12:00:00Z","endtime":"2020-01-01T12:00:00Z","description":"long dry spells"}]}]}`))
			})
		})

		When("everything is working and there are multiple cities", func() {
			It("should return a json weather forecast", func() {
				mockCache.EXPECT().Get("testcity").Return(nil, errors.New("cache miss"))
				mockCache.EXPECT().Get("testcity2").Return(nil, errors.New("cache miss"))
				mockReq, _ := http.NewRequest(http.MethodGet, "/weather?city=testcity,testcity2", nil)
				resp := httptest.NewRecorder()

				mockCoordinates.EXPECT().Find("testcity", "usa").Return(structs.CoOrdinates{}, nil)
				mockCoordinates.EXPECT().Find("testcity2", "usa").Return(structs.CoOrdinates{}, nil)

				setTime, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 12:00:00")
				weatherResult := structs.Weather{
					Start: setTime,
					End:   setTime,
				}
				weatherResult.Forecast.Long = "long dry spells"
				mockWeatherFetcher.EXPECT().Fetch(structs.CoOrdinates{}).Return([]structs.Weather{weatherResult}, nil).Times(2)

				mockCache.EXPECT().Store("testcity", gomock.Any())
				mockCache.EXPECT().Store("testcity2", gomock.Any())

				mockHandler.Handle(resp, mockReq)

				result := resp.Result()
				defer result.Body.Close()
				data, err := ioutil.ReadAll(result.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(data)).To(Equal(`{"forecast":[{"name":"testcity","detail":[{"starttime":"2020-01-01T12:00:00Z","endtime":"2020-01-01T12:00:00Z","description":"long dry spells"}]},{"name":"testcity2","detail":[{"starttime":"2020-01-01T12:00:00Z","endtime":"2020-01-01T12:00:00Z","description":"long dry spells"}]}]}`))
			})
		})

		When("everything is working", func() {
			It("should return a json weather forecast", func() {
				mockCache.EXPECT().Get("testcity").Return(nil, errors.New("cache miss"))
				mockReq, _ := http.NewRequest(http.MethodGet, "/weather?city=testcity", nil)
				resp := httptest.NewRecorder()

				mockCoordinates.EXPECT().Find("testcity", "usa").Return(structs.CoOrdinates{}, nil)

				setTime, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 12:00:00")
				weatherResult := structs.Weather{
					Start: setTime,
					End:   setTime,
				}
				weatherResult.Forecast.Long = "long dry spells"
				mockWeatherFetcher.EXPECT().Fetch(structs.CoOrdinates{}).Return([]structs.Weather{weatherResult}, nil)

				mockCache.EXPECT().Store("testcity", gomock.Any())

				mockHandler.Handle(resp, mockReq)

				result := resp.Result()
				defer result.Body.Close()
				data, err := ioutil.ReadAll(result.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(data)).To(Equal(`{"forecast":[{"name":"testcity","detail":[{"starttime":"2020-01-01T12:00:00Z","endtime":"2020-01-01T12:00:00Z","description":"long dry spells"}]}]}`))
			})
		})

		When("there is a cache hit", func() {
			It("should return from the cache and skip everything else", func() {
				pointInTime, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 15:00:00")
				mockCache.EXPECT().Get("testcity").Return([]handlerStructs.ResultForecast{
					handlerStructs.ResultForecast{
						Start:      pointInTime,
						End:        pointInTime,
						Prediction: "warm and sunny",
					},
				}, nil)

				mockReq, _ := http.NewRequest(http.MethodGet, "/weather?city=testcity", nil)
				resp := httptest.NewRecorder()

				mockHandler.Handle(resp, mockReq)

				result := resp.Result()
				defer result.Body.Close()
				data, err := ioutil.ReadAll(result.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(data)).To(Equal(`{"forecast":[{"name":"testcity","detail":[{"starttime":"2022-01-01T15:00:00Z","endtime":"2022-01-01T15:00:00Z","description":"warm and sunny"}]}]}`))
			})
		})
	})
})
