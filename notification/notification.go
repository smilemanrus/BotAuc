package notification

import (
	"BotAuc/storage"
)

type Processor struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Processor {
	return &Processor{
		storage: storage,
	}
}

func (p *Processor) Process(notyType string) error {
	//var err error
	//var actURL string
	//actualURLS := make(storage.UrlsAlias, 0)
	//for _, event := range events {
	//	log.Printf("got new event: %s", event.Text)
	//	actURL, err = p.SaveEvent(event)
	//	if err != nil {
	//		continue
	//	}
	//	actualURLS = append(actualURLS, actURL)
	//}
	//
	//if err = p.actualizeAucs(&actualURLS); err != nil {
	//	err = e.Wrap("can't process auc actualising", err)
	//}
	//
	return nil
}
