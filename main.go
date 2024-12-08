package main

import (
	"bufio"
	"fmt"
	"github.com/kiquetal/pokedexcli/internal"
	"os"
)

type ConfigUrl struct {
	Next     string
	Previous string
}
type cliCommand struct {
	name        string
	description string
	callback    func(url *ConfigUrl) error
}

func main() {

	config := ConfigUrl{}
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
	locations, next, previous, err := internal.GetLocations(c.Next)
	if err != nil {
		return err
	}

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
	locations, next, previous, err := internal.GetLocations(c.Previous)

	if err != nil {
		fmt.Printf("Error getting previous results: %s\n", err)
		return err
	}

	c.Next = next
	c.Previous = previous

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil

}
