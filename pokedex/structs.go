package main

import "github.com/Bybba/pokedex/internal/pokecache"

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, cache *pokecache.Cache, args []string) error
}

// map structs

type Config struct {
	Previous          string
	Next              string
	SupportedCommands map[string]cliCommand
	Pokedex           map[string]Pokemon
}

type LocationAreaResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type LocationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon Pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []Stat        `json:"stats"`
	Type           []PokemonType `json:"types"`
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Effor    int `json:"effort"`
	StatInfo struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}
