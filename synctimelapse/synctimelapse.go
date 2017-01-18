package main

import (
	"flag"
	"log"

	"github.com/dustin/octoprint"
)

func main() {
	base := flag.String("base", "", "octoprint base URL")
	token := flag.String("token", "", "octoprint token")

	flag.Parse()

	c, err := octoprint.New(*base, *token)
	if err != nil {
		log.Fatalf("Error setting up octoprint client: %v", err)
	}

	_, tls, err := c.ListTimelapses()
	if err != nil {
		log.Fatalf("Error listing timelapses: %v", err)
	}
	for _, tl := range tls {
		log.Printf(" %v (%v)", tl.Name, tl.SizeStr)
	}
}
