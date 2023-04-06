package main

import (
	eventConsumer "BotAuc/consumer/event-consumer"
	"BotAuc/events/auctions"
	"BotAuc/storage/sqlite"
	"context"
	"log"
)

const (
	tgHost      = "api.telegram.org"
	storagePath = "data/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("can't run db: %s", err)
	}
	if err = storage.Init(context.Background()); err != nil {
		log.Fatalf("can't init db: %s", err)
	}
	aucProcessor := auctions.New(storage)
	consumer := eventConsumer.New(aucProcessor, aucProcessor, 0, 10)
	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}
