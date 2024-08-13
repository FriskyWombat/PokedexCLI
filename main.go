package main

import (
	"fmt"
	"bufio"
	"os"
	"errors"
	"github.com/FriskyWombat/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name string
	description string
	callback func(cfg *config, client *pokeapi.Client) error
}
type config struct {
	nextLocUrl *string
	prevLocUrl *string
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
		"explore": {
			name: "explore", description: "Displays a list of Pokemon found at a location", callback: exploreCommand,
		},
	}

}

func helpCommand(cfg *config, client *pokeapi.Client) error {
	fmt.Println(" -----------------------------\n",
				 "Welcome to the Pokedex!\n",
				 "-----------------------------\n",
				 "Usage:")

	for _,cmd := range cmds {
		fmt.Println(" -" + cmd.name + ": " + cmd.description)
	}
	fmt.Println("\n -----------------------------")
	return nil
}

func exitCommand(cfg *config, client *pokeapi.Client) error {
	os.Exit(0)
	return nil
}


func mapCommand(cfg *config, client *pokeapi.Client) error {
	resp, err := client.FetchLocationData(*cfg.nextLocUrl)
	if err != nil {
		return err
	}
	for _, text := range resp.Results {
		fmt.Println(text.Name)
	}
	cfg.nextLocUrl = resp.Next
	cfg.prevLocUrl = resp.Previous
	return nil
}

func mapbCommand(cfg *config, client *pokeapi.Client) error {
	if cfg.prevLocUrl == nil {
		return errors.New("Err: No previous page exists")
	}
	resp, err := client.FetchLocationData(*cfg.prevLocUrl)
	if err != nil {
		return err
	}
	for _, text := range resp.Results {
		fmt.Println(text.Name)
	}
	cfg.nextLocUrl = resp.Next
	cfg.prevLocUrl = resp.Previous
	return nil
}

func exploreCommand(cfg *config, client *pokeapi.Client) error {
	os.Exit(0)
	return nil
}

func interpretCommand(text string, cfg *config, client *pokeapi.Client) error {
	if text == "" {
		return nil
	}
	for cmd := range cmds {
		if text == cmd {
			return cmds[cmd].callback(cfg, client)
		}
	}
	return errors.New("Err: Unknown command - " + text)
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	startUrl := pokeapi.GetFirstLocationUrl()
	cfg := config{&startUrl, nil}
	client := pokeapi.NewClient()
	for true {
		fmt.Printf("PKDX >:")
		in.Scan()
		err := interpretCommand(in.Text(), &cfg, &client)
		if err != nil {
			fmt.Println(err)
		}
	}
}