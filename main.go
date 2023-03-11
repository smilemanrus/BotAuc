package main

import (
	eventConsumer "BotAuc/consumer/event-consumer"
	"BotAuc/events/auctions"
	"log"
)

const (
	tgHost      = "api.telegram.org"
	storagePath = "storage"
	bachSize    = 100
)

func main() {
	//Бот
	//InitParams := initiation.InitiateParams()

	//tgClient := telegramClient.New(tgHost, InitParams.Token)
	//storage := files.New(storagePath)

	//eventsProcessor := tgProcessor.New(tgClient, storage)
	//
	//log.Print("service started")
	//consumer := eventConsumer.New(eventsProcessor, eventsProcessor, bachSize)
	//if err := consumer.Start(); err != nil {
	//	log.Fatal()
	//}

	//Парсер
	aucProcessor := auctions.New()
	consumer := eventConsumer.New(aucProcessor, aucProcessor, 0)
	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}
