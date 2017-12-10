package main

import (
	"github.com/BurntSushi/toml"
	"fmt"
	"log"
	//"bufio"
	"os"
)

type Dependencies []Dependency

type cf struct {
	Host    string       `toml:"host"`
	Imports Dependencies `toml:"import"`
}

type Dependency struct {
	Name       string `toml:"package"`
	Reference  string `toml:"version"`
	Repository string `toml:"repo"`
	VcsType    string `toml:"vcs"`
}

var c cf

func main() {

	if _, err := toml.DecodeFile("./ggopm.toml", &c); err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.Host)
	for _, s := range c.Imports {
		fmt.Printf("%s (%s)\n", s.Name, s.VcsType)
	}
	f, err := os.Create("./test.toml")
	if err !=nil {
		log.Fatal(err)
	}
	defer f.Close()

	e := toml.NewEncoder(f)
	if err = e.Encode(c); err!=nil{
		log.Fatal(err)
	}

}
