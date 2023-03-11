package data_consumer

import (
	"BotAuc/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	bathSize  int
}

func New(fetcher events.Fetcher, processor events.Processor, bathSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		bathSize:  bathSize,
	}
}
func (c Consumer) Start(cww string) error {
	for {
		gotAucs, err := c.fetcher.Fetch(c.bathSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s ", err.Error())
			continue
		}
		if len(gotAucs) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.SaveAucData(gotAucs); err != nil {
			log.Print(err.Error())
			continue
		}
	}
}
func (c Consumer) SaveAucData(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new auc: %s", event.Text)
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle auc: %s", err.Error())
			continue
		}
	}
	return nil
}
