package notification

import (
	"BotAuc/clients/telegram"
	"BotAuc/lib/e"
	"BotAuc/storage"
	"context"
	"errors"
	"strings"
)

type Processor struct {
	storage   storage.Storage
	messenger *telegram.Client
}

func New(storage storage.Storage, messenger *telegram.Client) *Processor {
	return &Processor{
		storage:   storage,
		messenger: messenger,
	}
}

func (p *Processor) Process(notyType string) error {

	switch notyType {
	case HourBeforeAuc:
		return p.alertAboutAucBfrHour(notyType)
	default:
		return errors.New("unknownAlert")
	}
}

func (p *Processor) alertAboutAucBfrHour(notyType string) error {
	ftrAucs, err := p.storage.GetAucsBfrHour(context.Background(), notyType)
	if err != nil {
		return e.Wrap("can't get aucs before hour", err)
	}
	for chatID, eventData := range ftrAucs {
		aucsMsg := strings.Join(eventData.Messages, "\n")
		if err = p.messenger.SendMessage(chatID, aucsMsg); err != nil {
			return e.Wrap("can't send aucs before hour", err)
		}
	}
	if err := p.storage.FixSendingAlert(context.Background(), ftrAucs, notyType); err != nil {
		return e.Wrap("can't fix sending alert", err)
	}
	return nil
}
