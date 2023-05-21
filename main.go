package main

import (
	telegramClient "BotAuc/clients/telegram"
	eventConsumer "BotAuc/consumer/event-consumer"
	"BotAuc/events/auctions"
	tgProcessor "BotAuc/events/telegram"
	"BotAuc/initiation"
	"BotAuc/notification"
	"BotAuc/storage/sqlite"
	"context"
	"log"
)

const (
	tgHost      = "api.telegram.org"
	storagePath = "data/sqlite/storage.db"
	bachSize    = 100
)

func main() {

	InitParams := initiation.InitiateParams()
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("can't run db: %s", err)
	}
	if err = storage.Init(context.Background()); err != nil {
		log.Fatalf("can't init db: %s", err)
	}

	tgClient := telegramClient.New(tgHost, InitParams.Token)

	//Бот
	eventsProcessor := tgProcessor.New(tgClient, storage)

	log.Print("service started")
	tgConsumer := eventConsumer.New(eventsProcessor, eventsProcessor, bachSize, 1)
	if err := tgConsumer.Start(); err != nil {
		log.Fatal()
	}

	//Парсер
	aucProcessor := auctions.New(storage)
	aucConsumer := eventConsumer.New(aucProcessor, aucProcessor, 0, 100)
	if err := aucConsumer.Start(); err != nil {
		log.Fatal()
	}

	//Оповещения
	notyProcessor := notification.New(storage, tgClient)
	if err := notyProcessor.Process(notification.HourBeforeAuc); err != nil {
		log.Fatal()
	}

}
