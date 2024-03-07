package main

import (
	"flag"
	"log"

	"taskbot/cmd/web/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	// tgClient
	tgClient := telegram.New(tgBotHost, mustToken())

	// fetcher

	// processor

	// consumer
}

func mustToken() (string) {
	token := flag.String(
		"token-bot-token", 
		"", 
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified.")
	}

	return token
}
