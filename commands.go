package main

import (
	"fmt"
	"math/rand"
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
		"catch": {
			name:        "catch",
			description: "Catches the given pokemon by name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Gives info of catched pokemon by name",
			callback:    commandInspect,
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
	fmt.Println("Exploring " + locationAreaDetails.Location.Name + "...")
	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationAreaDetails.PokemonEncounters {
		fmt.Println(" - " + pokemon.Pokemon.Name)
	}
	return nil
}

var pokedex = make(map[string]pokeapi.Pokemon)

func commandCatch(config *config, params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("no location name passed")
	}
	name := params[0]
	pokemon, err := pokeapi.GetPokemonInfo(name)
	if err != nil {
		return fmt.Errorf("failed executing cmd explore: %v", err)
	}
	fmt.Println("Throwing a Pokeball at " + pokemon.Name)
	rand := rand.Intn(1000)
	catched := pokemon.BaseExperience+rand < 500
	if catched {
		pokedex[pokemon.Name] = *pokemon
		fmt.Println(pokemon.Name + " was caught!")
	} else {
		fmt.Println(pokemon.Name + " escaped!")
	}
	return nil
}

func commandInspect(config *config, params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("no pokemon name passed")
	}
	pokemon, ok := pokedex[params[0]]
	if !ok {
		return fmt.Errorf("pokemon hasn't been catched yet")
	}
	fmt.Println("Name: " + pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, poketype := range pokemon.Types {
		fmt.Printf("  -%s\n", poketype.Type.Name)
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
