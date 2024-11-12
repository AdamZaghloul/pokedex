package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func main() {
	commands = getCommands()
	const prompt = "pokedex >"
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)

	for scanner.Scan() {

		command := scanner.Text()
		callback, err := parseCommand(command)
		if err != nil {
			fmt.Println(`Invalid command. Type "help" for list of commands`)
		} else {
			err = callback()
			if err != nil {
				break
			}
		}

		fmt.Print(prompt)
	}
}

func getCommands() map[string]cliCommand {
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

func parseCommand(command string) (func() error, error) {
	cmd, ok := commands[command]

	if ok {
		return cmd.callback, nil
	}
	return nil, errors.New("No such command")
}

func commandHelp() error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:\n")

	for _, command := range commands {
		fmt.Println(fmt.Sprintf("%s: %s", command.name, command.description))
	}

	fmt.Println("")

	return nil
}

func commandExit() error {
	return errors.New("Break from loop")
}
