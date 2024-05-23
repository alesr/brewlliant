package main

import (
	"log"

	"github.com/alesr/brewlliant/brewlliant"
)

func main() {
	if err := brewlliant.CheckBrew(); err != nil {
		log.Fatal(err)
	}

	if err := brewlliant.BrewList(); err != nil {
		log.Fatal(err)
	}

	if err := brewlliant.InstallFromBrewList(); err != nil {
		log.Fatal(err)
	}
}
