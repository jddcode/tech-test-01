package main

import (
	handlerWeather "github.com/jddcode/tech-test-ennismore/internal/handler-weather"
	"github.com/jddcode/tech-test-ennismore/internal/handler-weather/cache"
	"net/http"
)

func main() {
	cityCache := cache.New()
	http.HandleFunc("/weather", handlerWeather.New(cityCache).Handle)
	http.ListenAndServe(":8080", nil)
}
