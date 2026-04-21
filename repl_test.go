package main

import (
	"testing"
	"time"

	"github.com/Dorfieeee/bootdev-pokedex/internal/pokeapi"
	"github.com/Dorfieeee/bootdev-pokedex/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "--hello world",
			expected: []string{"--hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Test failed\nInput: '%v'\nActual: %v\nExpected: %v", c.input, actual, c.expected)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Test failed\nInput: '%v'\nActual: %v\nExpected: %v", c.input, actual, c.expected)
			}

		}
	}
}

func TestCommandInput(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "exit",
			expected: true,
		},
		{
			input:    "help",
			expected: true,
		},
		{
			input:    "explore",
			expected: true,
		},
		{
			input:    "badcommand",
			expected: false,
		},
	}

	for _, c := range cases {
		_, actual := getCommands()[c.input]
		if actual != c.expected {
			t.Errorf("Test failed\nInput: '%v'\nActual: %v\nExpected: %v", c.input, actual, c.expected)
		}
	}
}

func TestCommandInputWithArgs(t *testing.T) {
	cache := pokecache.NewCache(60 * time.Second)
	config := config{
		Api: pokeapi.NewPokeApi(cache),
	}
	exploreCmd, _ := getCommands()["explore"]
	err := exploreCmd.callback(&config)
	if err == nil {
		t.Errorf("Test failed: Expected error when calling `explore` command with no arguments")
	}
	if err.Error() != "Error: Expected location area name" {
		t.Errorf("Test failed: Expected error message 'Error: Expected location area name', got: %v", err)
	}
}
