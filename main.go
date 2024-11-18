package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal/pokecache"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
	Cache    pokecache.Cache
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
		Cache:    *pokecache.NewCache(60 * time.Second),
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
