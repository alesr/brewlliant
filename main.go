package main

import (
	"log"

	"github.com/alesr/rebrew/rebrew"
)

func main() {
	if err := rebrew.CheckBrew(); err != nil {
		log.Fatal(err)
	}

	if err := rebrew.BrewList(); err != nil {
		log.Fatal(err)
	}

	if err := rebrew.InstallFromBrewList(); err != nil {
		log.Fatal(err)
	}
}
