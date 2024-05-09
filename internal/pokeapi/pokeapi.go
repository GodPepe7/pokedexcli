package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/godpepe7/pokedexcli/internal/pokecache"
)

type LocationArea struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetails struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var cache = pokecache.NewCache(5 * time.Minute)

func fetch(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read reponse body")
	}
	return body, nil
}

func GetLocationAreas(url string) (*LocationArea, error) {
	body, inCache := cache.Get(url)
	locationArea := &LocationArea{}
	if !inCache {
		var err error
		body, err = fetch(url)
		if err != nil {
			return nil, err
		}
	}
	err := json.Unmarshal(body, &locationArea)
	if err != nil {
		return locationArea, fmt.Errorf("failed to unmarshal response body")
	}
	cache.Add(url, body)
	return locationArea, nil
}

func GetLocationAreaDetails(name string) (*LocationAreaDetails, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + name
	body, inCache := cache.Get(url)
	locationAreaDetails := &LocationAreaDetails{}
	if !inCache {
		var err error
		body, err = fetch(url)
		if err != nil {
			return nil, err
		}
	}
	err := json.Unmarshal(body, &locationAreaDetails)
	if err != nil {
		return locationAreaDetails, fmt.Errorf("failed to unmarshal response body")
	}
	cache.Add(url, body)
	return locationAreaDetails, nil
}

func GetPokemonInfo(name string) (*Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	body, inCache := cache.Get(url)
	pokemon := &Pokemon{}
	if !inCache {
		var err error
		body, err = fetch(url)
		if err != nil {
			return nil, err
		}
	}
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return pokemon, fmt.Errorf("failed to unmarshal response body")
	}
	cache.Add(url, body)
	return pokemon, nil
}
