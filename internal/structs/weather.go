package structs

import "time"

type Weather struct {
	Start, End time.Time
	IsDay bool
	TemperatureFarenheit int
	Wind struct {
		MinSpeed, MaxSpeed int
		Direction string
	}
	Forecast struct {
		Short, Long string
	}
}

func (w Weather) GetForecast() string {
	if len(w.Forecast.Long) > 0 {
		return w.Forecast.Long
	}
	return w.Forecast.Short
}
