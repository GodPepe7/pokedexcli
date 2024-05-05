package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("All commands:\n\n")
	commands := getAllCmds()
	for _, cliCommand := range commands {
		fmt.Println(cliCommand.name + ":" + cliCommand.description)
	}
	fmt.Print("\n")
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func ParseCommand(input string) error {
	commands := getAllCmds()
	cmd, ok := commands[input]
	if !ok {
		return fmt.Errorf("command doesn't exist")
	}
	err := cmd.callback()
	if err != nil {
		return fmt.Errorf("error executing command: %v", err)
	}
	return nil
}
