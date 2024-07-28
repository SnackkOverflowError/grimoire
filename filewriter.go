package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Grimoire struct {
	Spells map[string]Spell `json:"spells"`
}

func ReadGrimoireFromDisk() Grimoire {

	p := os.ExpandEnv("$HOME/.grimoire.json")
	g := Grimoire{}
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		g.Spells = make(map[string]Spell)
		return g
	}

	file, _ := os.ReadFile(p)
	json.Unmarshal(file, &g)

	if g.Spells == nil {
		g.Spells = make(map[string]Spell)
	}

	return g
}

func (g Grimoire) FlushToFile() {
	bytes, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		// what to do with this error?
		panic(err)
	}

	p := os.ExpandEnv("$HOME/.grimoire.json")

	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(p)
		if err != nil {
			panic(err)
		}
	}

	err = os.WriteFile(p, bytes, 0777)
	if err != nil {
		panic(err)
	}

}

func (g Grimoire) AddSpell(name string, spell Spell) {
	g.Spells[name] = spell
}

func (g Grimoire) GetSpell(name string) Spell {
	return g.Spells[name]
}
