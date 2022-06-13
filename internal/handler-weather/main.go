package handlerWeather

import (
	coOrdinateFinder "github.com/jddcode/tech-test-ennismore/internal/co-ordinate-finder"
	weatherFetcher "github.com/jddcode/tech-test-ennismore/internal/weather-fetcher"
)

func New(cache Cache) Handler {
	return handler{
		coOrdinates: coOrdinateFinder.New(),
		weather:     weatherFetcher.New(),
		cache:       cache,
	}
}
