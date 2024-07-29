package main

import (
	"log"

	server "github.com/alexPavlikov/go-wiki/cmd"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
