package coOrdinateFinder

import httpClient "github.com/jddcode/tech-test-ennismore/internal/http-client"

func New() Finder {
	return finder{
		web: httpClient.New(),
	}
}
