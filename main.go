package main

import (
	"log"

	"github.com/ekucher/telegram-bot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
