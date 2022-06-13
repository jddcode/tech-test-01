package httpClient

import (
	"io/ioutil"
	"net/http"
)

//go:generate mockgen -destination=../mocks/mock-http-client.go -package=mocks . Client
type Client interface {
	Get(url string) (string, error)
}

type client struct {}

func (c client) Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
