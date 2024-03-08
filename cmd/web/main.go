package main

import (
	"flag"
	"log"

	tgClient "taskbot/cmd/web/clients/telegram"
	 "taskbot/cmd/web/events/telegram"
	"taskbot/cmd/web/storage/mysql"
	"taskbot/cmd/web/consumer/event-consumer"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	dsn, token := Init()

	s, err := mysql.New(dsn)
	if err != nil {
		log.Fatal("Can't connect to storage:", err)
	}

	if err := s.Init(); err != nil {
		log.Fatal("can't init storage:", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, token),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}


}

func Init() (string, string) {
	dsn := flag.String(
		"dsn", 
		"", 
		"Data source name to access database.",
	)

	token := flag.String(
		"token", 
		"", 
		"token for access to telegram bot",
	)

	flag.Parse()

	if *dsn == "" {
		log.Fatal("DSN is not specified.")
	}

	if *token == "" {
		log.Fatal("Token is not specified.")
	}


	return *dsn, *token
}

