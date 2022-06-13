package cache

import "github.com/jddcode/tech-test-ennismore/internal/handler-weather/structs"

func New() Cache {
	return &cache{
		content: make(map[string][]structs.ResultForecast),
	}
}
