package main

import (
	"log"

	"github.com/eloylp/aton/internal/detector"
)

func main() {
	server, err := detector.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
