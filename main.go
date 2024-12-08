package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {

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
		err := command.callback()
		if err != nil {
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
	}

	return commands
}

func helpCommand() error {
	fmt.Printf("Welcome to the Pokedex !\n help: Display a help message\n exit: Exit the program\n")
	return nil
}

func exitCommand() error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}
