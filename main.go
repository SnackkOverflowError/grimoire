package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

/*
 * grimoire collect [name] [description] [function]
 *
 * This should open the prompt if you have an incorrect number of params
 *
 *
 *
 * grimoire cast [name]
 *
 * This should open a nice fzf list prompt if name doesnt exist or is not provided,
 * if name has been provided it should show the name, desc, and code for the "spell"
 *
 */

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("you need arguments")
		os.Exit(1)
	}

	switch args[0] {
	case "collect":
		collect_option(args[1:])
		break
	case "cast":
		fmt.Println("cast")
		cast_option(args[1:])
		break
	default:
		fmt.Println("argument unrecognized:", args[0])
	}

	os.Exit(0)
}

func cast_option(args []string) {
	grimoire := ReadGrimoireFromDisk()
	if len(args) == 0 {
		// NOTE no args, open the correct view
	}

	spell := grimoire.GetSpell(args[0])

	cmd, err := spell.Command(args[1:])

	if err != nil {
		log.Fatal("There was an Error constructing the Command:", err)
	}

	log.Println("Command constructed:", err)

	CastSpell(cmd)
}

func collect_option(args []string) {
	g := ReadGrimoireFromDisk()
	g.FlushToFile()

	if _, err := tea.NewProgram(initialModel(g)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

}
