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
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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
