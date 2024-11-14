package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
}

type mapResult struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

var commands map[string]cliCommand

func main() {
	config := Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

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
			err = callback(&config)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
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
	}
}

func parseCommand(command string) (func(*Config) error, error) {
	cmd, ok := commands[command]

	if ok {
		return cmd.callback, nil
	}
	return nil, errors.New("no such command")
}

func commandHelp(config *Config) error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:\n")

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println("")

	return nil
}

func commandExit(config *Config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *Config) error {
	body, err := httpGet(config.Next)
	if err != nil {
		return err
	}

	result := mapResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err

	}

	for _, item := range result.Results {
		fmt.Println(item.Name)
	}

	config.Next = result.Next
	config.Previous = result.Previous

	return nil
}

func commandMapb(config *Config) error {
	if config.Previous == "" {
		return errors.New(`already at the beginning of the map`)
	}
	body, err := httpGet(config.Previous)
	if err != nil {
		return err
	}

	result := mapResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err

	}

	for _, item := range result.Results {
		fmt.Println(item.Name)
	}

	config.Next = result.Next
	config.Previous = result.Previous

	return nil
}

func httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body))
	}
	if err != nil {
		return nil, err
	}

	return body, nil
}
