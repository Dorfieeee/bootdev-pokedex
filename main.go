package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	config := config{}
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
		cmd, ok := Commands[input[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := cmd.callback(&config)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
