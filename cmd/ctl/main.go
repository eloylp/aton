package main

import (
	"log"
	"os"

	"github.com/eloylp/aton/components/ctl"
)

func main() {
	c, err := ctl.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := c.AddMJPEGCapturer("capt1", os.Getenv("CAPT_URL"), 10); err != nil {
		log.Fatal(err)
	}
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
}
