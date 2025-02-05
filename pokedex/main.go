package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Bybba/pokedex/internal/pokecache"
)

// cleaning input -> lowercase, trim whitespace

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func main() {
	// Initialize shared cache and config once
	cache := pokecache.NewCache(10 * time.Second)
	config := &Config{
		SupportedCommands: make(map[string]cliCommand),
		Pokedex:           make(map[string]Pokemon),
	}

	rand.Seed(time.Now().UnixNano())
	// Create supported commands dynamically and updating the config with the generated commands
	supported_commands := createSupportedCommands()
	config.SupportedCommands = supported_commands

	// Start CLI loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Println("\nExiting...")
			return
		}
		input := scanner.Text()

		clean_input := cleanInput(input)

		// Handle empty input
		if len(clean_input) == 0 {
			fmt.Println("Please enter a command.")
			continue
		}

		// Look up the command in supported_commands
		cmd, exists := supported_commands[clean_input[0]]
		if !exists {
			fmt.Println("Unknown command.")
		} else {
			args := clean_input[1:]
			err := cmd.callback(config, cache, args)
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
			}
		}
	}
}
