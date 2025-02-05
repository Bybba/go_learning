package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Bybba/pokedex/internal/pokecache"
)

var fetchFromAPI func(url string) ([]byte, error) = realFetchFromAPI

func realFetchFromAPI(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("body read error: %s", err)
	}

	return body, nil
}

func getCachedOrFetchArea(url string, cache *pokecache.Cache) (LocationAreaResponse, error) {
	var location LocationAreaResponse

	// Check the cache first
	cachedData, exists := cache.Get(url)
	if exists {
		err := json.Unmarshal(cachedData, &location)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("unmarshal error from cache: %s", err)
		}
		return location, nil
	}

	// Fetch from API if not cached
	body, err := fetchFromAPI(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	// Parse into LocationAreaResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("unmarshal error: %s", err)
	}

	// Save the fetched result into cache
	cachedLocation, err := json.Marshal(location)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("marshal error: %s", err)
	}
	cache.Add(url, cachedLocation)

	return location, nil
}

func getCachedorFetchEncounters(location_area string, cache *pokecache.Cache) (LocationAreaDetail, error) {
	var encounters LocationAreaDetail
	url := "https://pokeapi.co/api/v2/location-area/" + location_area

	// check the cache first
	cachedData, exists := cache.Get(url)
	if exists {
		err := json.Unmarshal(cachedData, &encounters)
		if err != nil {
			return LocationAreaDetail{}, fmt.Errorf("unmarshal error from cache: %s", err)
		}
		return encounters, nil
	}

	body, err := fetchFromAPI(url)
	if err != nil {
		return LocationAreaDetail{}, err
	}

	err = json.Unmarshal(body, &encounters)
	if err != nil {
		return LocationAreaDetail{}, fmt.Errorf("unmarshal error: %s", err)
	}

	cached_encounter, err := json.Marshal(encounters)
	if err != nil {
		return LocationAreaDetail{}, fmt.Errorf("marshal error: %s", err)
	}

	cache.Add(url, cached_encounter)

	return encounters, nil
}
