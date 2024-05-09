package main

import (
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/godpepe7/pokedexcli/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, params ...string) error
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
			name:        "mapb",
			description: "Gives the previous 20 location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Gives all pokemons in given location",
			callback:    commandExplore,
		},
	}
}

func commandHelp(config *config, params ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("All commands:\n\n")
	commands := getAllCmds()
	for _, cliCommand := range commands {
		fmt.Println(cliCommand.name + ":" + cliCommand.description)
	}
	fmt.Print("\n")
	return nil
}

func commandExit(config *config, params ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config, params ...string) error {
	locationArea, err := pokeapi.GetLocationAreas(config.Next)
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

func commandMapB(config *config, params ...string) error {
	if config.Previous == "" {
		fmt.Println("No previous location areas")
		return nil
	}
	locationArea, err := pokeapi.GetLocationAreas(config.Previous)
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

func commandExplore(config *config, params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("no location name passed")
	}
	name := params[0]
	locationAreaDetails, err := pokeapi.GetLocationAreaDetails(name)
	if err != nil {
		return fmt.Errorf("failed executing cmd explore: %v", err)
	}
	fmt.Println("Exploring " + name + "...")
	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationAreaDetails.PokemonEncounters {
		fmt.Println(" - " + pokemon.Pokemon.Name)
	}
	return nil
}

func ParseCommand(input string) error {
	splitInput := strings.Split(input, " ")
	if len(splitInput) == 0 {
		return nil
	}
	commands := getAllCmds()
	cmd, ok := commands[splitInput[0]]
	if !ok {
		return fmt.Errorf("command doesn't exist")
	}
	args := splitInput[1:]
	err := cmd.callback(&cfg, args...)
	if err != nil {
		return fmt.Errorf("error executing command: %v", err)
	}
	return nil
}
