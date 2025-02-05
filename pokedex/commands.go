package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/Bybba/pokedex/internal/pokecache"
)

func commandExit(config *Config, cache *pokecache.Cache, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, cache *pokecache.Cache, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range config.SupportedCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(config *Config, cache *pokecache.Cache, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"

	if config.Next != "" {
		url = config.Next
	}

	location, err := getCachedOrFetchArea(url, cache)
	if err != nil {
		return fmt.Errorf("error fetching data: %s\n", err)
	}

	config.Next = location.Next
	config.Previous = location.Previous

	for _, loc := range location.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(config *Config, cache *pokecache.Cache, args []string) error {
	// Handle edge case when there is no previous page
	if config.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	location, err := getCachedOrFetchArea(config.Previous, cache)
	if err != nil {
		return fmt.Errorf("commandMapb error: %s", err)
	}

	config.Next = location.Next
	config.Previous = location.Previous

	for _, loc := range location.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandExplore(config *Config, cache *pokecache.Cache, args []string) error {
	if len(args) == 0 {
		fmt.Println("Please provide an area name in addition to the explore command")
		return nil
	}

	pokemons, err := getCachedorFetchEncounters(args[0], cache)
	if err != nil {
		return err
	}

	fmt.Println("Exploring " + args[0])
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemons.PokemonEncounters {
		fmt.Println("- " + pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config, cache *pokecache.Cache, args []string) error {
	if len(args) == 0 {
		fmt.Println("Provide a pokemon name with the catch command")
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	body, err := fetchFromAPI(url)
	if err != nil {
		return fmt.Errorf("could not fetch the data: %s\n", err)
	}

	var pokemon Pokemon

	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return fmt.Errorf("could not unmarshal: %s\n", err)
	}

	catch_try := rand.Intn(500)
	difficutly_level := pokemon.BaseExperience

	fmt.Println("Throwing a Pokeball at " + pokemon.Name + "...")

	if catch_try >= difficutly_level {
		fmt.Println(pokemon.Name + " was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
		config.Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Println(pokemon.Name + " escaped!")
	}

	return nil
}

func commandInspect(config *Config, cache *pokecache.Cache, args []string) error {
	if len(args) == 0 {
		fmt.Println("Please input a Pokemon name with the inspect command")
		return nil
	}

	pokemon_name := args[0]

	fmt.Printf("Debugging: %+v\n", config.Pokedex[pokemon_name])

	if pokemon, exists := config.Pokedex[pokemon_name]; exists {
		fmt.Println("Name: " + pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("-%s: %v\n", stat.StatInfo.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, poke_type := range pokemon.Type {
			fmt.Printf("- %s\n", poke_type.Type.Name)
		}
	} else {
		fmt.Println("You have not caught this Pokemon yet!")
	}

	return nil
}

func commandPokedex(config *Config, cache *pokecache.Cache, args []string) error {
	if len(config.Pokedex) == 0 {
		fmt.Println("You have not caught any Pokemon, user the catch")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.Pokedex {
		fmt.Println("- " + pokemon.Name)
	}
	return nil
}

func createSupportedCommands() map[string]cliCommand {
	// Declare supported_commands as a placeholder
	supported_commands := make(map[string]cliCommand)

	// Now initialize the commands using the placeholder
	supported_commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	supported_commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	supported_commands["map"] = cliCommand{
		name:        "map",
		description: "Lists 20 locations at a time",
		callback:    commandMap,
	}
	supported_commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Go back 20 locations if not on the first page",
		callback:    commandMapb,
	}
	supported_commands["explore"] = cliCommand{
		name:        "explore",
		description: "List all possible Pokemon encounters in an area",
		callback:    commandExplore,
	}
	supported_commands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempt to catch a pokemon",
		callback:    commandCatch,
	}
	supported_commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Provides info about a caught pokemon",
		callback:    commandInspect,
	}
	supported_commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "Lists all Pokemon in your Pokedex",
		callback:    commandPokedex,
	}

	return supported_commands
}
