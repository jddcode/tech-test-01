package structs

type ResultCity struct {
	City        string           `json:"name"`
	Predictions []ResultForecast `json:"detail"`
}
