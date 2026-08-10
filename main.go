package main

import (
	"log"
	"os"

	"github.com/ekucher/kbot/cmd"
)

func main() {
	if os.Getenv("TELE_TOKEN") == "" {
		log.Fatal("environment variable TELE_TOKEN is not set")
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
