package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type locationArea struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetMapAreas(url string) (locationArea, error) {
	locationArea := locationArea{}
	res, err := http.Get(url)
	if err != nil {
		return locationArea, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return locationArea, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return locationArea, fmt.Errorf("failed to read reponse body")
	}
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return locationArea, fmt.Errorf("failed to unmarshal response body")
	}
	return locationArea, nil
}
