package main

import (
	"flag"
	tgClient "goTelegram/client/telegram"
	event_consumer "goTelegram/consuner/event-consumer"
	"goTelegram/events/telegram"
	"goTelegram/storage/files"
	"log"
)

const (
	emptyToken  = ""
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(
		eventsProcessor,
		eventsProcessor,
		batchSize,
	)

	if err := consumer.Start(); err != nil {
		log.Fatal("error running consumer: ", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token",
	)

	flag.Parse()
	if *token == emptyToken {
		log.Fatal("empty token")
	}
	return *token
}
