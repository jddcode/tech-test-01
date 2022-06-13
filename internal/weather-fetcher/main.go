package weatherFetcher

import httpClient "github.com/jddcode/tech-test-ennismore/internal/http-client"

func New() WeatherFetcher {
	return weatherFetcher{
		web: httpClient.New(),
	}
}
