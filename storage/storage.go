package storage

import (
	"context"
	"time"
)

type Auction struct {
	Name      string
	StartDate time.Time
	EndDate   time.Time
	URL       string
	Status    string
}

type EventData struct {
	Message string
	URL     string
}

type EventsData map[int][]EventData

type UrlsAlias []string

type Storage interface {
	SaveData(ctx context.Context, p *Auction) error
	IsExists(ctx context.Context, p *Auction) (bool, error)
	ActualizeAucs(ctx context.Context, urls *UrlsAlias) error
	GetFutureAucs(ctx context.Context) (string, error)
	SubscrToAucs(ctx context.Context, chatID int, username string) error
	UnSubscrFormAucs(ctx context.Context, chatID int) error
	GetAucsBfrHour(ctx context.Context, eventType string) (EventsData, error)
	FixSendingAlert(ctx context.Context, urls EventsData, notyType string) error
}

func NewEventData(messages, url string) EventData {
	return EventData{
		Message: messages,
		URL:     url,
	}
}
