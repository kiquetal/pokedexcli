package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kiquetal/pokedexcli/internal"
	"github.com/kiquetal/pokedexcli/internal/pokecache"
	"github.com/mtslzr/pokeapi-go/structs"
	"os"
	"strings"
	"time"
)

type ConfigUrl struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
}
type cliCommand struct {
	name        string
	description string
	callback    func(url *ConfigUrl, args ...string) error
}

func main() {

	config := ConfigUrl{
		cache: pokecache.NewCache(50 * time.Second),
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Display the prompt
		fmt.Print("pokedexcli> ")
		//using newScanner
		scanner.Scan()
		text := scanner.Text()
		parts := strings.SplitN(text, " ", 2)
		commandName := parts[0]
		var args []string
		if len(parts) > 1 {
			args = strings.Split(parts[1], " ")
		}
		allCommands := getCommands()
		command, ok := allCommands[commandName]
		if !ok {
			fmt.Println("Command not found")
			continue
		}
		err := command.callback(&config, args...)
		if err != nil {
			fmt.Printf("Error executing command: %s\n", err)
			fmt.Println("Error executing command")
		}

	}

}

func getCommands() map[string]cliCommand {

	commands := map[string]cliCommand{
		"help": cliCommand{
			name:        "help",
			description: "Display a help help message",
			callback:    helpCommand,
		},
		"exit": cliCommand{
			name:        "exit",
			description: "Exit the program",
			callback:    exitCommand,
		},
		"map": cliCommand{
			name:        "map",
			description: "Return location ares",
			callback:    mapCommand,
		},
		"mapb": cliCommand{
			name:        "mapb",
			description: "Return previous location areas",
			callback:    mapBCommand,
		},
		"cache": cliCommand{
			name:        "cache",
			description: "Display cache",
			callback:    cacheCommand,
		},
		"explore": cliCommand{
			name:        "explore",
			description: "Explore a location",
			callback:    exploreCommand,
		},
		"catch": cliCommand{
			name:        "catch",
			description: "Catch a pokemon",
			callback:    catchCommand,
		},
	}

	return commands
}

func catchCommand(url *ConfigUrl, args ...string) error {
	return nil
}

func cacheCommand(url *ConfigUrl, args ...string) error {
	fmt.Printf("previos value: %s\n", url.Previous)
	fmt.Printf("next value: %s\n", url.Next)
	//print all the keys from the cache
	keys := url.cache.Cache
	for k := range keys {
		fmt.Println(k)
	}
	return nil
}

func helpCommand(c *ConfigUrl, args ...string) error {
	fmt.Printf("Welcome to the Pokedex!\n help: Display a help message\n exit: Exit the program\n")
	return nil
}

func exitCommand(c *ConfigUrl, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func mapCommand(c *ConfigUrl, args ...string) error {
	fmt.Println("Getting location areas")

	// If Next is empty, this is the first request
	if c.Next == "" {
		locations, next, previous, err := internal.GetLocations("")
		if err != nil {
			fmt.Printf("Error getting locations: %s\n", err)
			return err
		}

		// Marshal and cache the results
		data, err := json.Marshal(locations)
		if err != nil {
			fmt.Println("Error marshalling locations")
			return err
		}

		fmt.Printf("Adding to cache key:%s", c.Next)
		c.cache.Add("https://pokeapi.co/api/v2/location-area?offset=0&limit=20", data)
		fmt.Println("First request")

		c.Next = next
		c.Previous = previous

		for _, location := range locations {
			fmt.Println(location.Name)
		}
		return nil
	}

	// Check cache for subsequent requests
	cachedValue, found := c.cache.Get(c.Next)
	if found {
		fmt.Println("Using cached value")
		var locations []structs.Result
		if err := json.Unmarshal(cachedValue, &locations); err != nil {
			fmt.Println("Error unmarshalling cached value")
			return err
		}
		for _, location := range locations {
			fmt.Println(location.Name)
		}
		return nil
	}

	// If not in cache, get from API
	locations, next, previous, err := internal.GetLocations(c.Next)
	if err != nil {
		fmt.Printf("Error getting locations: %s\n", err)
		return err
	}

	// Marshal and cache the results
	data, err := json.Marshal(locations)
	if err != nil {
		fmt.Println("Error marshalling locations")
		return err
	}

	c.cache.Add(c.Next, data)
	fmt.Printf("Adding to cache key:%s\n", c.Next)
	c.Next = next
	c.Previous = previous

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func mapBCommand(c *ConfigUrl, args ...string) error {
	fmt.Println("Getting the Previous Results")

	if c.Previous == "" {
		fmt.Println("No previous results")
		return nil
	}

	// Check cache first
	fmt.Printf("Looking for key:%s\n", c.Previous)
	cachedValue, found := c.cache.Get(c.Previous)
	if found {
		fmt.Println("Using cached value")
		fmt.Printf("Previous key:%s\n", c.Previous)
		var locations []structs.Result
		if err := json.Unmarshal(cachedValue, &locations); err != nil {
			fmt.Println("Error unmarshalling cached value")
			return err
		}
		for _, location := range locations {
			fmt.Println(location.Name)
		}
		return nil
	}

	// If not in cache, get from API
	locations, next, previous, err := internal.GetLocations(c.Previous)
	if err != nil {
		fmt.Printf("Error getting previous results: %s\n", err)
		return err
	}

	// Marshal data before caching
	data, err := json.Marshal(locations)
	if err != nil {
		fmt.Println("Error marshalling locations")
		return err
	}

	// Cache the results and update URLs

	fmt.Printf("Adding to cache key:%s\n", c.Previous)
	c.cache.Add(c.Previous, data)
	c.Next = next
	c.Previous = previous

	// Display locations
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil

}

func exploreCommand(c *ConfigUrl, args ...string) error {
	fmt.Println("Exploring a location")
	location := args[0]

	//check if the location is in the cache
	cachedValue, found := c.cache.Get(location)
	if found {
		fmt.Println("Using cached value")
		fmt.Printf("Location:%s\n", location)
		var pokemons []internal.Pokemon
		if err := json.Unmarshal(cachedValue, &pokemons); err != nil {
			fmt.Println("Error unmarshalling cached value")
			return err
		}
		fmt.Printf("Location Area:%s\n", location)
		return nil
	}

	// If not in cache, get from API
	pokemonsFound, err := internal.GetLocationArea(location)
	if err != nil {
		fmt.Printf("Error getting location area: %s\n", err)
		return err
	}

	// Marshal data before caching
	data, err := json.Marshal(pokemonsFound)

	if err != nil {
		fmt.Println("Error marshalling locations")
		return err
	}

	// Cache the results

	c.cache.Add(location, data)
	fmt.Println("Found Pokemon:")
	for _, encounter := range pokemonsFound {
		fmt.Printf("- %s\n", encounter.Name)
	}

	return nil
}
