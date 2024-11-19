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
	callback    func(*Config, string) error
}

type Config struct {
	Next     string
	Previous string
	Explore  string
	Pokemon  string
	Pokedex  map[string]Pokemon
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
		Explore:  "https://pokeapi.co/api/v2/location-area/",
		Pokemon:  "https://pokeapi.co/api/v2/pokemon/",
		Pokedex:  map[string]Pokemon{},
		Cache:    *pokecache.NewCache(60 * time.Second),
	}

	commands = getCommands()
	const prompt = "pokedex >"
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)

	for scanner.Scan() {

		command := scanner.Text()
		callback, args, err := parseCommand(command)
		if err != nil {
			fmt.Println(`Invalid command. Type "help" for list of commands`)
		} else {
			err = callback(&config, args)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}

		fmt.Print(prompt)
	}
}
