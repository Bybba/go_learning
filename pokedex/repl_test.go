package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Bybba/pokedex/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input  string
		output []string
	}{
		{
			input:  "   hello    world    ",
			output: []string{"hello", "world"},
		}, {
			input:  "1 2   3    ",
			output: []string{"1", "2", "3"},
		}, {
			input:  "",
			output: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.output) {
			t.Fatalf("output length not matching with expected: got %d, expect %d", len(actual), len(c.output))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.output[i]
			if word != expectedWord {
				t.Errorf("%s did not match %s", word, expectedWord)
			} else if i == len(actual)-1 {
				t.Logf("test passed: input %v matches %v", actual, c.output)
			}
		}
	}
}

func TestCommandMapWithCacheFound(t *testing.T) {
	args := []string{}
	
	cache := pokecache.NewCache(5 * time.Second)

	url := "https://pokeapi.co/api/v2/location-area/"

	cachedResponse := LocationAreaResponse{
		Results:  []Location{{Name: "test-location"}},
		Next:     "next_url",
		Previous: "previous_url",
	}

	cachedData, _ := json.Marshal(cachedResponse)
	cache.Add(url, cachedData)

	config := &Config{Next: url}

	err := commandMap(config, cache, args)
	if err != nil {
		t.Fatalf("commandMap failed: %s", err)
	}

	if config.Next != "next_url" {
		t.Errorf("expected config.Next to be 'next_url' but got '%s'", config.Next)
	}

	if config.Previous != "previous_url" {
		t.Errorf("expected config.Previous to be 'previous_url' but got '%s'", config.Previous)
	}
}

func TestCommandMapWithCacheMissing(t *testing.T) {
	args := []string{}
	
	cache := pokecache.NewCache(5 * time.Second)

	url := "https://pokeapi.co/api/v2/location-area/"
	config := &Config{Next: url}

	fetchFromAPI = func(url string) ([]byte, error) {
		apiResponse := LocationAreaResponse{
			Results:  []Location{{Name: "new_location"}},
			Next:     "next_url",
			Previous: "previous_url",
		}
		return json.Marshal(apiResponse)
	}

	err := commandMap(config, cache, args)
	if err != nil {
		t.Fatalf("commandMap failed: %s", err)
	}

	if config.Next != "next_url" {
		t.Errorf("expected config.Next to be 'next_url' but got '%s'", config.Next)
	}

	if config.Previous != "previous_url" {
		t.Errorf("expected config.Previous to be 'previous_url' but got '%s'", config.Previous)
	}

	cachedData, ok := cache.Get(url)
	if !ok {
		t.Fatalf("expected data to be cached, but it was not found")
	}

	var cachedResponse LocationAreaResponse
	err = json.Unmarshal(cachedData, &cachedResponse)
	if err != nil {
		t.Fatalf("failed to unmarshal cached data: %s", err)
	}

	if len(cachedResponse.Results) != 1 || cachedResponse.Results[0].Name != "new_location" {
		t.Errorf("expected cached result to contain 'new_location', but got '%v'", cachedResponse.Results)
	}

	// Ensure the cache is holding the correct "next" and "previous" URLs
	if cachedResponse.Next != "next_url" {
		t.Errorf("expected cached response Next to be 'next_url' but got '%s'", cachedResponse.Next)
	}

	if cachedResponse.Previous != "previous_url" {
		t.Errorf("expected cached response Previous to be 'previous_url' but got '%s'", cachedResponse.Previous)
	}
}

func TestCacheReapLoop(t *testing.T) {
	const interval = 2 * time.Second
	cache := pokecache.NewCache(interval)

	// Add an entry to the cache
	url := "https://pokeapi.co/api/v2/test-endpoint"
	data := []byte("test data")
	cache.Add(url, data)

	// Confirm entry exists initially
	_, ok := cache.Get(url)
	if !ok {
		t.Fatalf("expected entry to exist immediately after adding")
	}

	// Wait for longer than the interval, so the reap loop can clean it up
	time.Sleep(3 * time.Second)

	// Confirm entry no longer exists
	_, ok = cache.Get(url)
	if ok {
		t.Errorf("entry should have been removed after interval but was found")
	}
}
