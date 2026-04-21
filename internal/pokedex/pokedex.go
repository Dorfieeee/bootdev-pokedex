package pokedex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Dorfieeee/bootdev-pokedex/internal/pokeapi"
)

type Store struct {
	Pokemons map[string]pokeapi.Pokemon
}

type Pokedex struct {
	Store    Store
	Filepath string
}

func NewPokedex() *Pokedex {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := &Pokedex{
		Filepath: filepath.Join(cwd, "data"),
		Store: Store{
			Pokemons: make(map[string]pokeapi.Pokemon),
		},
	}
	if err := p.load(); err != nil {
		fmt.Println("Unable to load Pokedex Store")
	}
	return p
}

func (p *Pokedex) GetPokemon(name string) (pokeapi.Pokemon, bool) {
	pokemon, exists := p.Store.Pokemons[name]
	if !exists {
		return pokeapi.Pokemon{}, false
	}
	return pokemon, true
}

func (p *Pokedex) AddPokemon(pokemon pokeapi.Pokemon) {
	p.Store.Pokemons[pokemon.Name] = pokemon
	err := p.save()
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Pokedex) load() error {
	_, err := os.Stat(p.getStorePath())
	if err != nil || os.IsNotExist(err) {
		// first load or no data saved yet
		os.Mkdir(p.Filepath, 0755)
		return nil
	}
	raw, err := os.ReadFile(p.getStorePath())
	if err != nil {
		return err
	}
	var store Store
	if err := json.Unmarshal(raw, &store); err != nil {
		return err
	}
	p.Store = store
	return nil
}

func (p *Pokedex) save() error {
	raw, err := json.Marshal(p.Store)
	if err != nil {
		return err
	}
	if err := os.WriteFile(p.getTempStorePath(), raw, 0644); err != nil {
		return err
	}
	if err := os.Rename(p.getTempStorePath(), p.getStorePath()); err != nil {
		return err
	}
	return nil
}

func (p *Pokedex) getStorePath() string {
	return filepath.Join(p.Filepath, "store.json")
}

func (p *Pokedex) getTempStorePath() string {
	return filepath.Join(p.Filepath, "store.temp")
}
