package main

import (
	"log"

	"github.com/eloylp/aton/components/node"
)

func main() {
	server, err := node.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
