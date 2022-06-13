package coOrdinateFinder

import (
	"encoding/json"
	"errors"
	"fmt"
	httpClient "github.com/jddcode/tech-test-ennismore/internal/http-client"
	"github.com/jddcode/tech-test-ennismore/internal/structs"
	"net/url"
	"strconv"
)

const (
	ErrorNoCity       = "You must supply a city"
	ErrorNoCountry    = "You must supply a country"
	ErrorHTTPGet      = "HTTP GET error: %s"
	ErrorUnmarshall   = "Unmarshal error: %s"
	ErrorNoData       = "No data found after unmarshal"
	ErrorBadLatitude  = "Unrecognised latitude: %s"
	ErrorBadLongitude = "Unrecognised longitude: %s"
)

//go:generate mockgen -destination=../mocks/mock-co-ordinate-finder.go -package=mocks . Finder
type Finder interface {
	Find(city, country string) (structs.CoOrdinates, error)
}

type finder struct {
	web httpClient.Client
}

func (f finder) Find(city, country string) (structs.CoOrdinates, error) {
	if len(city) < 1 {
		return structs.CoOrdinates{}, errors.New(ErrorNoCity)
	}

	if len(country) < 1 {
		return structs.CoOrdinates{}, errors.New(ErrorNoCountry)
	}

	res, err := f.web.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s,%s&format=json", url.QueryEscape(city), url.QueryEscape(country)))
	if err != nil {
		return structs.CoOrdinates{}, fmt.Errorf(ErrorHTTPGet, err.Error())
	}

	data := result{}
	if err = json.Unmarshal([]byte(res), &data); err != nil {
		return structs.CoOrdinates{}, fmt.Errorf(ErrorUnmarshall, err.Error())
	}

	if len(data) < 1 {
		return structs.CoOrdinates{}, errors.New(ErrorNoData)
	}

	myLat, err := strconv.ParseFloat(data[0].Lat, 64)
	if err != nil {
		return structs.CoOrdinates{}, fmt.Errorf(ErrorBadLatitude, data[0].Lat)
	}

	myLon, err := strconv.ParseFloat(data[0].Lon, 64)
	if err != nil {
		return structs.CoOrdinates{}, fmt.Errorf(ErrorBadLongitude, data[0].Lon)
	}

	return structs.CoOrdinates{
		Latitude:  myLat,
		Longitude: myLon,
	}, nil
}
