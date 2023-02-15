package main

import (
	"flag"
	"log"
)

const (
	tgHost = "api.telegram.org"
)

func main() {
	//tgClient := telegram.New(tgHost, mustToken())

	//aucFetcher

	//processor

	//consumer.Start(fetcher, processor)
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
