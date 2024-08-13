package main

import (
	"fmt"
	"bufio"
	"os"
	"errors"
	"strings"
	"math/rand"
	"embed"
	"image"
	"github.com/FriskyWombat/pokedex/internal/pokeapi"
	"github.com/dolmen-go/kittyimg"
  _ "image/png"
)

//go:embed "35.png"
var files embed.FS

type cliCommand struct {
	name string
	description string
	callback func(args []string, cfg *config, client *pokeapi.Client) error
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
		"catch": {
			name: "catch", description: "Attempt to catch a pokemon", callback: catchCommand,
		},
	}
}

func helpCommand(args []string, cfg *config, client *pokeapi.Client) error {
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
		return errors.New("Err: No previous page exists")
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
	if len(resp.PokemonEncounters) == 0{
		return fmt.Errorf("This location has no Pokemon") 
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
	fmt.Println("Throwing a Pokeball at", pkmn.Name)
	if rand.Intn(pkmn.BaseExperience) < (15 + pkmn.BaseExperience/6)  {
		fmt.Println("You caught " + pkmn.Name + "!")
		/* Display an image of the caught Pokemon */
		f, err := files.Open("35.png")
		if err != nil {
			return err
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		kittyimg.Fprintln(os.Stdout, img)
	} else {
		fmt.Println("Shoot! It was so close too...")
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
	for true {
		fmt.Printf("PKDX >:")
		in.Scan()
		err := interpretCommand(in.Text(), &cfg, &client)
		if err != nil {
			fmt.Println(err)
		}
	}
}