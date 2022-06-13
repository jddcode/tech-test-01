package handlerWeather

import (
	"encoding/json"
	"fmt"
	coOrdinateFinder "github.com/jddcode/tech-test-ennismore/internal/co-ordinate-finder"
	"github.com/jddcode/tech-test-ennismore/internal/handler-weather/structs"
	weatherFetcher "github.com/jddcode/tech-test-ennismore/internal/weather-fetcher"
	"net/http"
	"strings"
	"time"
)

const (
	ErrorNoCities      = "Please supply a comma delimited list of cities as the URL parameter 'city'"
	ErrorNoCoordinates = "Could not find co-ordinates for city: %s"
	ErrorNoForecast    = "Could not get a weather forecast for the city: %s"
	ErrorMashallResult = "Could not marshall result into valid json: %s"
)

type Cache interface {
	Get(city string) ([]structs.ResultForecast, error)
	Store(city string, prediction []structs.ResultForecast)
}

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	coOrdinates coOrdinateFinder.Finder
	weather     weatherFetcher.WeatherFetcher
	cache       Cache
}

func (h handler) Handle(w http.ResponseWriter, r *http.Request) {
	cities := strings.Split(r.URL.Query().Get("city"), ",")
	if len(cities) < 1 || len(cities[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrorNoCities))
		return
	}

	output := structs.Result{}
	for _, city := range cities {
		if data, err := h.cache.Get(city); err == nil {
			output.Data = append(output.Data, structs.ResultCity{
				City:        city,
				Predictions: data,
			})
			continue
		}

		pos, err := h.coOrdinates.Find(city, "usa")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(ErrorNoCoordinates, city)))
			return
		}

		forecasts, err := h.weather.Fetch(pos)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(ErrorNoForecast, city)))
			return
		}

		predictions := make([]structs.ResultForecast, 0)
		for _, forecast := range forecasts {
			if forecast.Start.After(time.Now().Add(time.Hour * 48)) {
				break
			}

			predictions = append(predictions, structs.ResultForecast{
				Start:      forecast.Start,
				End:        forecast.End,
				Prediction: forecast.GetForecast(),
			})
		}

		h.cache.Store(city, predictions)
		output.Data = append(output.Data, structs.ResultCity{
			City:        city,
			Predictions: predictions,
		})
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(ErrorMashallResult, err.Error())))
		return
	}

	w.Write(bytes)
}
