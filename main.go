package main

import (
	telegramClient "BotAuc/clients/telegram"
	eventConsumer "BotAuc/consumer/event-consumer"
	tgProcessor "BotAuc/events/telegram"
	"BotAuc/initiation"
	"BotAuc/storage/files"
	"log"
)

const (
	tgHost      = "api.telegram.org"
	storagePath = "storage"
	bachSize    = 100
)

func main() {

	InitParams := initiation.InitiateParams()

	tgClient := telegramClient.New(tgHost, InitParams.Token)
	storage := files.New(storagePath)

	eventsProcessor := tgProcessor.New(tgClient, storage)

	log.Print("service started")
	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, bachSize)
	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

}
