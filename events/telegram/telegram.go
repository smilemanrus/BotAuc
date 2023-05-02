package telegram

import (
	"BotAuc/clients/telegram"
	eventsLib "BotAuc/events"
	"BotAuc/lib/e"
	"BotAuc/storage"
	"errors"
)

type Processor struct {
	tg             *telegram.Client
	offset         int
	storage        storage.Storage
	ProcessMessage interface{}
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

func New(client *telegram.Client, s storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: s,
	}
}

func (p *Processor) Fetch(limit int) ([]eventsLib.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]eventsLib.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(events []eventsLib.Event) error {
	var err error
	for _, event := range events {
		if err = processEvent(p, event); err != nil {
			err = e.Wrap("can't process event", err)
			break
		}
	}
	return err
}

func processEvent(p *Processor, event eventsLib.Event) error {
	switch event.Type {
	case eventsLib.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("Can't process message", ErrUnknownEventType)
	}
}

func event(upd telegram.Update) eventsLib.Event {
	updType := fetchType(upd)
	res := eventsLib.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == eventsLib.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) eventsLib.Type {
	if upd.Message == nil {
		return eventsLib.Unknown
	}
	return eventsLib.Message
}

func (p *Processor) processMessage(event eventsLib.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	if err := p.doCMD(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil
}

func meta(event eventsLib.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", eventsLib.ErrUnknownMetaType())
	}
	return res, nil
}
