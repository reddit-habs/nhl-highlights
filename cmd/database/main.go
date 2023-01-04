package main

import (
	"log"
	"os"

	"github.com/sbstp/nhl-highlights/repository"
)

func main() {
	dbpath := os.Args[1]
	repo, err := repository.New(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	if err = repo.Close(); err != nil {
		log.Fatal(err)
	}

}
