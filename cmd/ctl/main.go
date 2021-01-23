package main

import (
	"log"

	"github.com/eloylp/aton/components/ctl"
)

func main() {
	c, err := ctl.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
}
