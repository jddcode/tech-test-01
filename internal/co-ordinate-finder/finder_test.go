package coOrdinateFinder

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jddcode/tech-test-ennismore/internal/mocks"
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
		mockFinder     finder
	)

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
		mockHttpClient = mocks.NewMockClient(mockController)
		mockFinder = finder{
			web: mockHttpClient,
		}
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Context("Fetching the co-ordinates for a city", func() {
		When("the length of the city is zero", func() {
			It("should return an error", func() {
				_, err := mockFinder.Find("", "USA")
				Expect(err).To(Equal(errors.New(ErrorNoCity)))
			})
		})

		When("the length of the city is zero", func() {
			It("should return an error", func() {
				_, err := mockFinder.Find("New York", "")
				Expect(err).To(Equal(errors.New(ErrorNoCountry)))
			})
		})

		When("there is an error calling the web service", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return("", errors.New("error carrying out GET request"))
				_, err := mockFinder.Find("New York", "USA")
				Expect(err).To(Equal(fmt.Errorf(ErrorHTTPGet, "error carrying out GET request")))
			})
		})

		When("there is an error unmarshalling the response from the web service", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return("---", nil)
				_, err := mockFinder.Find("New York", "USA")
				Expect(err).To(Equal(fmt.Errorf(ErrorUnmarshall, "invalid character '-' in numeric literal")))
			})
		})

		When("the data returned from the web service cannot be unmarshalled into usable information", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return("[]", nil)
				_, err := mockFinder.Find("New York", "USA")
				Expect(err).To(Equal(errors.New(ErrorNoData)))
			})
		})

		When("the data returned from the web service has a latitude which does not properly convert to a float64", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`[{"lat":"invalid-lat"}]`, nil)
				_, err := mockFinder.Find("New York", "USA")
				Expect(err).To(Equal(fmt.Errorf(ErrorBadLatitude, "invalid-lat")))
			})
		})

		When("the data returned from the web service has a longitude which does not properly convert to a float64", func() {
			It("should return an error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`[{"lat":"1.23", "lon":"invalid-lon"}]`, nil)
				_, err := mockFinder.Find("New York", "USA")
				Expect(err).To(Equal(fmt.Errorf(ErrorBadLongitude, "invalid-lon")))
			})
		})

		When("the data can be unmarshalled and makes sense, and the lat and long are valid", func() {
			It("should return the co-ordinates with no error", func() {
				mockHttpClient.EXPECT().Get(gomock.Any()).Return(`[{"lat":"1.23", "lon":"1.23"}]`, nil)
				pos, err := mockFinder.Find("New York", "USA")
				Expect(err).ToNot(HaveOccurred())
				Expect(pos.Longitude).To(Equal(1.23))
				Expect(pos.Latitude).To(Equal(1.23))
			})
		})
	})
})
