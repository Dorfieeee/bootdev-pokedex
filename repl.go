package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/Dorfieeee/bootdev-pokedex/internal/pokeapi"
	"github.com/Dorfieeee/bootdev-pokedex/internal/pokedex"
)

type config struct {
	Next     *string
	Previous *string
	Api      *pokeapi.PokeApi
	Pokedex  *pokedex.Pokedex
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args ...string) error
	usage       string
}

func cleanInput(text string) []string {
	var input []string
	input = strings.Fields(strings.ToLower(text))
	return input
}

func commandExit(c *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(url *string, c *config) error {
	pgList, err := c.Api.GetLocationAreas(url)
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

func commandMapNext(c *config, args ...string) error {
	if c.Next == nil && c.Previous != nil {
		fmt.Println("you're on the last page")
	}
	return commandMap(c.Next, c)
}

func commandMapBack(c *config, args ...string) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return commandMap(c.Previous, c)
}

func commandExplore(c *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Error: Expected location area name")
	}
	name := args[0]
	fmt.Printf("Exploring %s...\n", name)
	data, err := c.Api.GetLocationArea(args[0])
	if err != nil {
		return fmt.Errorf("Error loading location area `%s`: %w", name, err)
	}
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Error: Expected Pokemon name")
	}
	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	pokemon, err := c.Api.GetPokemon(name)
	if err != nil {
		return fmt.Errorf("Error loading `%s` information: %w", name, err)
	}

	difficulty := (math.Log(float64(pokemon.BaseExperience)) / float64(pokemon.BaseExperience)) * 10
	chance := rand.Float64()

	if chance < difficulty {
		// Sucess
		fmt.Printf("%s was caught!\n", name)
		c.Pokedex.AddPokemon(pokemon)
	} else {
		fmt.Printf("%s escaped!\n", name)
	}

	return nil
}

func commandInspect(c *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Error: Expected Pokemon name")
	}
	name := args[0]
	pokemon, prs := c.Pokedex.GetPokemon(name)
	if !prs {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Printf("Stats:\n")
	for _, s := range pokemon.Stats {
		fmt.Printf("- %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("- %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(c *config, args ...string) error {
	if len(c.Pokedex.Store.Pokemons) == 0 {
		fmt.Println("No records in your Pokedex yet")
		return nil
	}
	fmt.Println("Your Pokedex")
	for name := range c.Pokedex.Store.Pokemons {
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations",
			callback:    commandMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explores the given area and displays all Pokemons to encounter",
			callback:    commandExplore,
			usage:       "Usage: explore <location-area>",
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a Pokemon",
			callback:    commandCatch,
			usage:       "Usage: catch <pokemon>",
		},
		"inspect": {
			name:        "inspect",
			description: "Shows basic information about a Pokemon in your Pokedex",
			callback:    commandInspect,
			usage:       "Usage: inspect <pokemon>",
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all Pokemon in your Pokedex",
			callback:    commandPokedex,
		},
	}
}
