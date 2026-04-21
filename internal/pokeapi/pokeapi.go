package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Dorfieeee/bootdev-pokedex/internal/pokecache"
)

const PokeApiBaseUrl = "https://pokeapi.co/api/v2"

const (
	LocationAreaRoute = "/location-area"
	PokemonRoute      = "/pokemon"
)

type PokeApi struct {
	Client *http.Client
	Cache  *pokecache.Cache
}

type ResourceReference struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreas struct {
	Count    int                 `json:"count"`
	Next     *string             `json:"next"`
	Previous *string             `json:"previous"`
	Results  []ResourceReference `json:"results"`
}

type LocationArea struct {
	Name              string            `json:"name"`
	Location          ResourceReference `json:"location"`
	PokemonEncounters []struct {
		Pokemon ResourceReference `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int               `json:"base_stat"`
		Stat     ResourceReference `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int               `json:"slot"`
		Type ResourceReference `json:"type"`
	} `json:"types"`
}

func NewPokeApi(cache *pokecache.Cache) *PokeApi {
	return &PokeApi{
		Client: &http.Client{},
		Cache:  cache,
	}
}

func (api *PokeApi) Get(url string) ([]byte, error) {
	if data, cached := api.Cache.Get(url); cached {
		return data, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Request error: %w", err)
	}
	res, err := api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Network error: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %s", res.Status)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read error: %w", err)
	}
	api.Cache.Add(url, data)

	return data, nil
}

func (api *PokeApi) GetLocationAreas(url *string) (LocationAreas, error) {
	baseUrl := fmt.Sprintf("%s%s", PokeApiBaseUrl, LocationAreaRoute)
	if url != nil {
		baseUrl = *url
	}
	raw, err := api.Get(baseUrl)
	if err != nil {
		return LocationAreas{}, err
	}

	var data LocationAreas
	if err := json.Unmarshal(raw, &data); err != nil {
		return LocationAreas{}, fmt.Errorf("Unmarshal error: %w", err)
	}

	return data, nil
}

func (api *PokeApi) GetLocationArea(name string) (LocationArea, error) {
	url := fmt.Sprintf("%s%s/%s", PokeApiBaseUrl, LocationAreaRoute, name)
	raw, err := api.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	var data LocationArea
	if err := json.Unmarshal(raw, &data); err != nil {
		return LocationArea{}, fmt.Errorf("Unmarshal error: %w", err)
	}

	return data, nil
}

func (api *PokeApi) GetPokemon(name string) (Pokemon, error) {
	url := fmt.Sprintf("%s%s/%s", PokeApiBaseUrl, PokemonRoute, name)
	raw, err := api.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	var data Pokemon
	if err := json.Unmarshal(raw, &data); err != nil {
		return Pokemon{}, fmt.Errorf("Unmarshal error: %w", err)
	}

	return data, nil
}
