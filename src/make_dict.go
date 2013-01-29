package main

import (
	darts "github.com/awsong/go-darts"
	"log"
)

func main() {
	_, err := darts.Import("data/deals.txt", "data/deals.lib", true)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(dict)

	log.Println("done")
}
