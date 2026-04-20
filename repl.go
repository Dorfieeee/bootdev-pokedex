package main

import (
	"fmt"
	"os"
	"strings"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

var Commands = make(map[string]cliCommand)

func cleanInput(text string) []string {
	var input []string
	input = strings.Fields(strings.ToLower(text))
	return input
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range Commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(url string, c *config) error {
	pgList, err := getPaginatedList[locationArea](url)
	if err != nil {
		return fmt.Errorf("Error loading location areas: %w", err)
	}
	c.Next = pgList.Next
	c.Previous = pgList.Previous
	for _, loc := range pgList.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapNext(c *config) error {
	if c.Next == "" && c.Previous != "" {
		fmt.Println("you're on the last page")
	}
	url := LocationAreaRoute
	if c.Next != "" {
		url = c.Next
	}
	return commandMap(url, c)
}

func commandMapBack(c *config) error {
	url := c.Previous
	if url == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return commandMap(url, c)
}

func init() {
	Commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	Commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	Commands["map"] = cliCommand{
		name:        "map",
		description: "Displays next 20 locations",
		callback:    commandMapNext,
	}
	Commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays previous 20 locations",
		callback:    commandMapBack,
	}
}
