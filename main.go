package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Dorfieeee/bootdev-pokedex/internal/pokeapi"
	"github.com/Dorfieeee/bootdev-pokedex/internal/pokecache"
	"github.com/Dorfieeee/bootdev-pokedex/internal/pokedex"
)

func main() {
	cache := pokecache.NewCache(60 * time.Second)
	config := config{
		Api:     pokeapi.NewPokeApi(cache),
		Pokedex: pokedex.NewPokedex(),
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if token := scanner.Scan(); !token {
			continue
		}
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		cmd, ok := getCommands()[input[0]]
		args := input[1:]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := cmd.callback(&config, args...)
		if err != nil {
			fmt.Println(err.Error())
			if cmd.usage != "" {
				fmt.Println(cmd.usage)
			}
		}
	}
}
