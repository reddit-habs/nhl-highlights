package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/sbstp/nhl-highlights/nhlapi"
)

func main() {
	for _, team := range nhlapi.Teams {
		fmt.Println(team)
	}

	client := nhlapi.NewClient()
	resp, err := client.Schedule("", "")
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
}
