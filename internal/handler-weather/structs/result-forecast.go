package structs

import "time"

type ResultForecast struct {
	Start      time.Time `json:"starttime"`
	End        time.Time `json:"endtime"`
	Prediction string    `json:"description"`
}
