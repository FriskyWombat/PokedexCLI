package main

import (
	"fmt"
	"bufio"
	"os"
	"errors"
	"net/http"
)

type cliCommand struct {
	name string
	description string
	callback func(cfg *config) error
}
type config struct {
	nextUrl *string
	prevUrl *string
}

var cmds map[string]cliCommand = make(map[string]cliCommand)

func init() {
	cmds = map[string]cliCommand {
		"help": {
			name: "help", description: "Displays a help message", callback:helpCommand,
		},
		"exit": {
			name: "exit", description: "Exit the Pokedex", callback: exitCommand,
		},
		"map": {
			name: "map", description: "Displays next 20 locations", callback: mapCommand,
		},
		"mapb": {
			name: "mapb", description: "Displays the previous 20 locations", callback: mapbCommand,
		},
	}

}

func helpCommand(cfg *config) error {
	fmt.Println(" -----------------------------\n",
				 "Welcome to the Pokedex!\n",
				 "-----------------------------\n",
				 "Usage:\n")

	for _,cmd := range cmds {
		fmt.Println(" -" + cmd.name + ": " + cmd.description)
	}
	fmt.Println("\n -----------------------------\n")
	return nil
}

func exitCommand(cfg *config) error {
	os.Exit(0)
	return nil
}

func fetchMapData(url string, cfg *config) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	data := FetchedData{}
	err = parseResponse(res, &data)
	if err != nil {
		return err
	}
	for _, r := range data.Results {
		fmt.Println(r.Name)
	}
	cfg.nextUrl = data.Next
	cfg.prevUrl = data.Previous
	return nil
}

func mapCommand(cfg *config) error {
	return fetchMapData(*cfg.nextUrl, cfg)
}

func mapbCommand(cfg *config) error {
	return fetchMapData(*cfg.prevUrl, cfg)
}

func interpretCommand(text string, cfg *config) error {
	for cmd := range cmds {
		if text == cmd {
			return cmds[cmd].callback(cfg)
		}
	}
	return errors.New("Err: Unknown command")
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	startURL := "https://pokeapi.co/api/v2/location?offset=0&limit=20"
	cfg := config{&startURL, nil}
	for true {
		fmt.Printf("PKDX >:")
		in.Scan()
		interpretCommand(in.Text(), &cfg)
	}
}