package main

import (
	"bufio"
	"errors"
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"strings"

	"github.com/FriskyWombat/pokedex/internal/pokeapi"
	"github.com/qeesung/image2ascii/convert"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string, cfg *config, client *pokeapi.Client) error
}
type config struct {
	nextLocUrl *string
	prevLocUrl *string
}

var cmds map[string]cliCommand = make(map[string]cliCommand)
var dex map[string]pokeapi.Pokemon = make(map[string]pokeapi.Pokemon)

func init() {
	cmds = map[string]cliCommand{
		"help": {
			name: "help", description: "Displays a help message", callback: helpCommand,
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
		"catch": {
			name: "catch", description: "Attempt to catch a pokemon", callback: catchCommand,
		},
		"inspect": {
			name: "inspect", description: "Display information about your caught Pokemon", callback: inspectCommand,
		},
		"pokedex": {
			name: "pokedex", description: "Display a list of all of the Pokemon you've caught", callback: pokedexCommand,
		},
	}
}

func helpCommand(args []string, cfg *config, client *pokeapi.Client) error {
	fmt.Println(" -----------------------------\n",
		"Welcome to the Pokedex!\n",
		"-----------------------------\n",
		"Usage:")

	for _, cmd := range cmds {
		fmt.Println(" -" + cmd.name + ": " + cmd.description)
	}
	fmt.Println("\n -----------------------------")
	return nil
}

func exitCommand(args []string, cfg *config, client *pokeapi.Client) error {
	os.Exit(0)
	return nil
}

func mapCommand(args []string, cfg *config, client *pokeapi.Client) error {
	resp, err := client.FetchLocationData(*cfg.nextLocUrl)
	if err != nil {
		return err
	}
	for _, text := range resp.Results {
		fmt.Println(" - " + text.Name)
	}
	cfg.nextLocUrl = resp.Next
	cfg.prevLocUrl = resp.Previous
	return nil
}

func mapbCommand(args []string, cfg *config, client *pokeapi.Client) error {
	if cfg.prevLocUrl == nil {
		return errors.New("err: No previous page exists")
	}
	resp, err := client.FetchLocationData(*cfg.prevLocUrl)
	if err != nil {
		return err
	}
	for _, text := range resp.Results {
		fmt.Println(" - " + text.Name)
	}
	cfg.nextLocUrl = resp.Next
	cfg.prevLocUrl = resp.Previous
	return nil
}

func exploreCommand(args []string, cfg *config, client *pokeapi.Client) error {
	resp, err := client.FetchLocationAreaData(args[1])
	if err != nil {
		return err
	}
	if len(resp.PokemonEncounters) == 0 {
		return fmt.Errorf("this location has no Pokemon")
	}
	for _, pkmn := range resp.PokemonEncounters {
		fmt.Println(" - " + pkmn.Pokemon.Name)
	}
	return nil
}

func catchCommand(args []string, cfg *config, client *pokeapi.Client) error {
	pkmn, err := client.FetchPokemonData(args[1])
	if err != nil {
		return err
	}
	fmt.Println("Throwing a Pokeball at", pkmn.Name+"...")
	if rand.Intn(pkmn.BaseExperience) < (25 + pkmn.BaseExperience/5) {
		fmt.Println("You caught " + pkmn.Name + "!")
		dex[pkmn.Name] = pkmn
		convertOptions := convert.DefaultOptions
		convertOptions.FixedHeight = 10
		convertOptions.FixedWidth = 20
		client.PrintPokemonImage(&pkmn, &convertOptions)
	} else {
		fmt.Println("Shoot! It was so close, too...")
	}
	return nil
}

func inspectCommand(args []string, cfg *config, client *pokeapi.Client) error {
	pkmn, ok := dex[args[1]]
	if !ok {
		return fmt.Errorf("no pokedex data found for %s - you need to catch one first", args[1])
	}
	convertOptions := convert.DefaultOptions
	convertOptions.FixedHeight = 30
	convertOptions.FixedWidth = 60
	client.PrintPokemonImage(&pkmn, &convertOptions)
	fmt.Println("=== " + strings.ToUpper(pkmn.Name) + " ===")
	if len(pkmn.Types) == 1 {
		fmt.Printf("Type: %s\n", pkmn.Types[0].Type.Name)
	} else {
		fmt.Printf("Types: %s/%s\n", pkmn.Types[0].Type.Name, pkmn.Types[1].Type.Name)
	}
	fmt.Printf("Height: %d   Weight: %d\n", pkmn.Height, pkmn.Weight)
	fmt.Printf("Base stats:\n -HP: %d\n -ATK: %d\n -DEF: %d\n -SPATK: %d\n -SPDEF: %d\n -SPD: %d\n",
		pkmn.Stats[0].BaseStat, pkmn.Stats[1].BaseStat, pkmn.Stats[2].BaseStat, pkmn.Stats[3].BaseStat, pkmn.Stats[4].BaseStat, pkmn.Stats[5].BaseStat)
	return nil
}

func pokedexCommand(args []string, cfg *config, client *pokeapi.Client) error {
	if len(dex) == 0 {
		return fmt.Errorf("you haven't caught any pokemon yet")
	}
	fmt.Println("Your Pokedex:")
	for name, pkmn := range dex {
		fmt.Println(" - " + name)
		convertOptions := convert.DefaultOptions
		convertOptions.FixedHeight = 5
		convertOptions.FixedWidth = 10
		client.PrintPokemonImage(&pkmn, &convertOptions)
	}
	return nil
}

func interpretCommand(text string, cfg *config, client *pokeapi.Client) error {
	if text == "" {
		return nil
	}
	args := strings.Fields(text)
	for cmd := range cmds {
		if args[0] == cmd {
			return cmds[cmd].callback(args, cfg, client)
		}
	}
	return errors.New("Err: Unknown command - " + text)
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	startUrl := pokeapi.GetFirstLocationUrl()
	cfg := config{&startUrl, nil}
	client := pokeapi.NewClient()
	for {
		fmt.Printf("PKDX >:")
		in.Scan()
		err := interpretCommand(in.Text(), &cfg, &client)
		if err != nil {
			fmt.Println(err)
		}
	}
}
