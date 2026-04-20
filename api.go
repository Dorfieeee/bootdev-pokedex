package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type paginatedList[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T
}

type locationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

const (
	LocationAreaRoute string = "https://pokeapi.co/api/v2/location-area"
)

func getPaginatedList[T any](url string) (paginatedList[T], error) {
	res, err := http.Get(url)
	if err != nil {
		return paginatedList[T]{}, fmt.Errorf("Network error while requesting %s", url)
	}

	defer res.Body.Close()

	var data paginatedList[T]
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return paginatedList[T]{}, fmt.Errorf("Unable to parse response body from %s", url)
	}

	return data, nil
}
