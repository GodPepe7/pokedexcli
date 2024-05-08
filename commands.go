package main

import (
	"fmt"
	"os"

	pokeapi "github.com/godpepe7/pokedexcli/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config) error
}

var cfg config = config{
	Next:     "https://pokeapi.co/api/v2/location-area/",
	Previous: "",
}

func getAllCmds() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas and advances with each call",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Gives the previous 20 location areas",
			callback:    commandMapB,
		},
	}
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("All commands:\n\n")
	commands := getAllCmds()
	for _, cliCommand := range commands {
		fmt.Println(cliCommand.name + ":" + cliCommand.description)
	}
	fmt.Print("\n")
	return nil
}

func commandExit(config *config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config) error {
	locationArea, err := pokeapi.GetMapAreas(config.Next)
	if err != nil {
		return err
	}
	for _, val := range locationArea.Results {
		fmt.Println(val.Name)
	}
	if locationArea.Previous == nil {
		config.Previous = ""
	} else {
		config.Previous = *locationArea.Previous
	}
	config.Next = locationArea.Next
	return nil
}

func commandMapB(config *config) error {
	if config.Previous == "" {
		fmt.Println("No previous location areas")
		return nil
	}
	locationArea, err := pokeapi.GetMapAreas(config.Previous)
	if err != nil {
		return err
	}
	for _, val := range locationArea.Results {
		fmt.Println(val.Name)
	}
	config.Next = locationArea.Next
	if locationArea.Previous == nil {
		config.Previous = ""
	} else {
		config.Previous = *locationArea.Previous
	}
	return nil
}

func ParseCommand(input string) error {
	commands := getAllCmds()
	cmd, ok := commands[input]
	if !ok {
		return fmt.Errorf("command doesn't exist")
	}
	err := cmd.callback(&cfg)
	if err != nil {
		return fmt.Errorf("error executing command: %v", err)
	}
	return nil
}
