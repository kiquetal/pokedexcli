package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kiquetal/pokedexcli/internal"
	"github.com/kiquetal/pokedexcli/internal/pokecache"
	"github.com/mtslzr/pokeapi-go/structs"
	"os"
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
	callback    func(url *ConfigUrl) error
}

func main() {

	config := ConfigUrl{
		cache: pokecache.NewCache(30 * time.Second),
	}

	for {
		// Display the prompt
		fmt.Print("pokedexcli> ")
		//using newScanner
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		allCommands := getCommands()
		command, ok := allCommands[text]
		if !ok {
			fmt.Println("Command not found")
			continue
		}
		err := command.callback(&config)
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
	}

	return commands
}

func helpCommand(c *ConfigUrl) error {
	fmt.Printf("Welcome to the Pokedex !\n help: Display a help message\n exit: Exit the program\n")
	return nil
}

func exitCommand(c *ConfigUrl) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func mapCommand(c *ConfigUrl) error {
	fmt.Println("Getting location areas")

	// Check cache first
	if c.Next != "" {
		cachedValue, found := c.cache.Get(c.Next)
		if found {
			fmt.Println("Using cached value")
			var locations []structs.Result
			e := json.Unmarshal(cachedValue, &locations)
			if e != nil {
				fmt.Println("Error unmarshalling cached value")
				return e
			}
			for _, location := range locations {
				fmt.Println(location.Name)
			}
			return nil
		}
	}

	// If not in cache, get from API
	locations, next, previous, err := internal.GetLocations(c.Next)
	if err != nil {
		fmt.Printf("Error getting locations: %s\n", err)
		return err
	}

	data, err2 := json.Marshal(locations)
	if err2 != nil {
		fmt.Println("Error marshalling locations")
		return err2
	}

	// Cache the current results
	c.cache.Add(c.Next, data)

	// Update the URLs after caching
	c.Next = next
	c.Previous = previous

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func mapBCommand(c *ConfigUrl) error {
	fmt.Println("Getting the Previous Results")

	// previous is an interface{} type, probably a Results struct
	if c.Previous == "" {
		fmt.Println("No previous results")
		return nil
	}

	cachedValue, found := c.cache.Get(c.Previous)
	fmt.Println("Cached value: ", cachedValue)
	if found {
		fmt.Println("Using cached value")
		var locations []structs.Result
		e := json.Unmarshal(cachedValue, &locations)
		if e != nil {
			fmt.Println("Error unmarshalling cached value")
			return e
		}
		for _, location := range locations {
			fmt.Println(location.Name)
		}
		return nil
	}

	locations, next, previous, err := internal.GetLocations(c.Previous)

	if err != nil {
		fmt.Printf("Error getting previous results: %s\n", err)
		return err
	}

	data, err2 := json.Marshal(locations)
	c.cache.Add(previous, data)

	c.Next = next
	c.Previous = previous
	if err2 != nil {
		fmt.Println("Error marshalling locations")
		return err2
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil

}
