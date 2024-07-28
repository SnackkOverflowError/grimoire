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

	g := Grimoire{}
	if _, err := os.Stat("~/.grimoire.json"); errors.Is(err, os.ErrNotExist) {
		return g
	}

	file, _ := os.ReadFile("~/.grimoire.json")
	json.Unmarshal(file, &g)

	return g
}

func (g Grimoire) FlushToFile() {

	bytes, err := json.Marshal(g)
	if err != nil {
		// what to do with this error?
		panic(err)
	}

	err = os.WriteFile("~/.grimoire.json", bytes, 0644)
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
