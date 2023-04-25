package event_consumer

import (
	"BotAuc/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher    events.Fetcher
	processor  events.Processor
	bathSize   int
	pauseValue time.Duration
}

func New(fetcher events.Fetcher, processor events.Processor, bathSize int, pauseValue time.Duration) Consumer {
	return Consumer{
		fetcher:    fetcher,
		processor:  processor,
		bathSize:   bathSize,
		pauseValue: pauseValue,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.bathSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s ", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(c.pauseValue * time.Second)
			continue
		}
		if err := c.HandleEvents(gotEvents); err != nil {
			log.Print(err.Error())
			continue
		}
	}
}

func (c Consumer) HandleEvents(events []events.Event) error {
	if err := c.processor.Process(events); err != nil {
		log.Printf("can't handle event: %s", err.Error())
	}
	return nil
}
