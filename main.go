package main

import (
	tgClient "BotAuc/clients/telegram"
	event_consumer "BotAuc/consumer/event-consumer"
	tgProcessor "BotAuc/events/telegram"
	"BotAuc/storage/files"
	"flag"
	"log"
)

const (
	tgHost      = "api.telegram.org"
	storagePath = "storage"
	bachSize    = 100
)

func main() {
	tgClient := tgClient.New(tgHost, mustToken())
	storage := files.New(storagePath)

	eventsProcessor := tgProcessor.New(tgClient, storage)

	log.Print("service started")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, bachSize)
	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

}

func mustToken() string {
	//botAuc -tgTokenBot 'my token'
	token := flag.String(
		"tgTokenBot",
		"",
		"Token for Tg send-bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("Token not found")
	}
	return *token
}
