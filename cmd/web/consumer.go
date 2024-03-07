package main

import (
	"flag"
	"log"
)

func main() {
	// tgClient
	t := mustToken()

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
