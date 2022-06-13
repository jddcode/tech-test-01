package weatherFetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	httpClient "github.com/jddcode/tech-test-ennismore/internal/http-client"
	"github.com/jddcode/tech-test-ennismore/internal/structs"
	fetcherStructs "github.com/jddcode/tech-test-ennismore/internal/weather-fetcher/structs"
	"strconv"
	"strings"
	"time"
)

const (
	ErrorGetRequest         = "Error fetching weather report via GET: %s"
	ErrorUnmarshalLookup    = "Error unmarshalling the co-ordinate weather lookup: %s"
	ErrorNoForecastResource = "Error finding the forecast resource from the co-ordinate weather lookup"
	ErrorGetForecast        = "Error fetching forecast via GET: %s"
	ErrorUnmarshalForecast  = "Error unarmshalling the forecast: %s"
	ErrorUnusualWindSpeed   = "Error converting wind speed measures to integers: %s"
	ErrorUnusualStartTime   = "Error converting start time to time.Time: %s"
	ErrorUnusualEndTime     = "Error converting end time to time.Time: %s"
)

//go:generate mockgen -destination=../mocks/mock-weather-fetcher.go -package=mocks . WeatherFetcher
type WeatherFetcher interface {
	Fetch(pos structs.CoOrdinates) ([]structs.Weather, error)
}

type weatherFetcher struct {
	web httpClient.Client
}

func (w weatherFetcher) Fetch(pos structs.CoOrdinates) ([]structs.Weather, error) {
	resp, err := w.web.Get(fmt.Sprintf("https://api.weather.gov/points/%.5f,%.5f", pos.Latitude, pos.Longitude))
	if err != nil {
		return nil, fmt.Errorf(ErrorGetRequest, err.Error())
	}

	lookupResult := fetcherStructs.ResponseCoOrdinateLookup{}
	err = json.Unmarshal([]byte(resp), &lookupResult)
	if err != nil {
		return nil, fmt.Errorf(ErrorUnmarshalLookup, err.Error())
	}

	if len(lookupResult.Properties.Forecast) < 1 {
		return nil, errors.New(ErrorNoForecastResource)
	}

	resp, err = w.web.Get(lookupResult.Properties.Forecast)
	if err != nil {
		return nil, fmt.Errorf(ErrorGetForecast, err.Error())
	}

	forecastData := fetcherStructs.ResponseForecast{}
	err = json.Unmarshal([]byte(resp), &forecastData)
	if err != nil {
		return nil, fmt.Errorf(ErrorUnmarshalForecast, err.Error())
	}

	myWeather := make([]structs.Weather, 0)
	for _, period := range forecastData.Properties.Periods {
		weather := structs.Weather{
			IsDay:                period.IsDaytime,
			TemperatureFarenheit: period.Temperature,
		}

		weather.Start, err = w.parseTimeString(period.StartTime)
		if err != nil {
			return nil, fmt.Errorf(ErrorUnusualStartTime, err.Error())
		}

		weather.End, err = w.parseTimeString(period.EndTime)
		if err != nil {
			return nil, fmt.Errorf(ErrorUnusualEndTime, err.Error())
		}

		weather.Wind.MinSpeed, weather.Wind.MaxSpeed, err = w.getWindSpeeds(period.WindSpeed)
		if err != nil {
			return nil, fmt.Errorf(ErrorUnusualWindSpeed, err.Error())
		}

		weather.Wind.Direction = period.WindDirection

		weather.Forecast.Short = period.ShortForecast
		weather.Forecast.Long = period.DetailedForecast
		myWeather = append(myWeather, weather)
	}
	return myWeather, nil
}

func (w weatherFetcher) getWindSpeeds(windSpeedStr string) (int, int, error) {
	parts := strings.Split(windSpeedStr, " ")
	switch len(parts) {
	case 2:
		minSpeed, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, errors.New("Could not convert min speed to integer")
		}
		return minSpeed, minSpeed, nil
	case 4:
		minSpeed, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, errors.New("Could not convert min speed to integer")
		}

		maxSpeed, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, 0, errors.New("Could not convert max speed to integer")
		}
		return minSpeed, maxSpeed, nil
	default:
		return 0, 0, errors.New("Unexpected format for wind speed string")
	}
}

func (w weatherFetcher) parseTimeString(timeString string) (time.Time, error) {
	if len(timeString) < 18 {
		return time.Time{}, errors.New("invalid time string")
	}
	myTime, err := time.Parse("2006-01-02T15:04:05", timeString[0:19])
	if err != nil {
		return time.Time{}, err
	}
	return myTime, nil
}
