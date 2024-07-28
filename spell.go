package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Spell struct {
	Description      string   `json:"desc"`
	Command_template string   `json:"templ"`
	Args             []string `json:"args"`
}

func NewSpell(desc string, cmd_tmpl string) (Spell, error) {

	spell := Spell{}
	spell.Description = desc

	spell.Command_template = cmd_tmpl
	spell.Args = extractArgs(cmd_tmpl)

	return spell, nil
}

const MATCH_ARGS_REG = "<<<([^<> ]*?)>>>"

func extractArgs(cmd_tmpl string) []string {
	r, _ := regexp.Compile(MATCH_ARGS_REG)
	return r.FindAllString(cmd_tmpl, -1)
}

func injectArgs(cmd_tmpl string, arg_names []string, arg_values []string) (string, error) {
	if len(arg_names) != len(arg_values) {
		return "", errors.New("argument list length mismatch")
	}

	var cmd string = cmd_tmpl

	for index, name := range arg_names {
		tmpl := fmt.Sprintf("<<<%s>>>", name)
		cmd = strings.ReplaceAll(cmd, tmpl, arg_values[index])
	}

	return cmd, nil
}

func (s Spell) Command(args []string) (string, error) {
	cmd, err := injectArgs(s.Command_template, s.Args, args)
	if err != nil {
		return "", err
	}

	cmd = strings.ReplaceAll(cmd, "\n", " ; ")

	return cmd, nil

}

func CastSpell(cmd_string string) error {
	cmd := exec.Command("bash", "-c", cmd_string)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatal("There was an error running the command string:", err)
	}

	return nil
}
