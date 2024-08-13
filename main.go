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
	callback func(cfg *config) error
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


func mapCommand(cfg *config) error {
	resp, err := pokeapi.FetchMapData(*cfg.nextLocUrl)
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

func mapbCommand(cfg *config) error {
	if cfg.prevLocUrl == nil {
		return errors.New("Err: No previous page exists")
	}
	resp, err := pokeapi.FetchMapData(*cfg.prevLocUrl)
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

func interpretCommand(text string, cfg *config) error {
	if text == "" {
		return nil
	}
	for cmd := range cmds {
		if text == cmd {
			return cmds[cmd].callback(cfg)
		}
	}
	return errors.New("Err: Unknown command - " + text)
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	startUrl := pokeapi.GetFirstLocationUrl()
	cfg := config{&startUrl, nil}
	for true {
		fmt.Printf("PKDX >:")
		in.Scan()
		err := interpretCommand(in.Text(), &cfg)
		if err != nil {
			fmt.Println(err)
		}
	}
}