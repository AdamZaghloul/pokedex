package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

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
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore AREA-NAME",
			description: "Displays the names of all pokemon in te given AREA-NAME argument",
			callback:    commandExplore,
		},
	}
}

func parseCommand(command string) (func(*Config, string) error, string, error) {
	args := strings.Split(command, " ")

	cmd, ok := commands[args[0]]

	if len(args) == 1 {
		args = append(args, "")
	}

	if ok {
		return cmd.callback, args[1], nil
	}
	return nil, "", errors.New("no such command")
}

func commandHelp(config *Config, arg string) error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:")
	fmt.Println()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println("")

	return nil
}

func commandExit(config *Config, arg string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *Config, arg string) error {
	var err error

	body, ok := config.Cache.Get(config.Next)
	if !ok {
		body, err = httpGet(config.Next)
		if err != nil {
			return err
		}

		err = config.Cache.Add(config.Next, body)
		if err != nil {
			return err
		}
	}

	result := mapResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err

	}

	fmt.Println()

	for _, item := range result.Results {
		fmt.Println(item.Name)
	}

	fmt.Println()

	config.Next = result.Next
	config.Previous = result.Previous

	return nil
}

func commandMapb(config *Config, arg string) error {
	var err error

	if config.Previous == "" {
		return errors.New(`already at the beginning of the map`)
	}
	body, ok := config.Cache.Get(config.Previous)
	if !ok {
		body, err = httpGet(config.Previous)
		if err != nil {
			return err
		}

		err = config.Cache.Add(config.Previous, body)
		if err != nil {
			return err
		}
	}

	result := mapResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err

	}

	fmt.Println()

	for _, item := range result.Results {
		fmt.Println(item.Name)
	}

	fmt.Println()

	config.Next = result.Next
	config.Previous = result.Previous

	return nil
}

func commandExplore(config *Config, arg string) error {
	var err error

	if arg == "" {
		return errors.New("no area specified.")
	}

	body, ok := config.Cache.Get(config.Explore + arg)

	if !ok {
		body, err = httpGet(config.Explore + arg)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return errors.New("invalid area name.")
			} else {
				return err
			}
		}

		err = config.Cache.Add(config.Explore+arg, body)
		if err != nil {
			return err
		}
	}

	result := Explore{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err

	}

	fmt.Println()
	fmt.Println("Exploring " + arg + "...")
	fmt.Println("Found Pokemon:")
	fmt.Println()

	for _, encounter := range result.PokemonEncounters {
		fmt.Println(" - " + encounter.Pokemon.Name)
	}

	fmt.Println()

	return nil
}
