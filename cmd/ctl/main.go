package main

import (
	"fmt"
	"os"

	"github.com/eloylp/aton/internal/ctl"
)

func main() {
	c, err := ctl.NewFromEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := c.AddMJPEGCapturer("capt1", os.Getenv("CAPT_URL"), 10); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := c.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
