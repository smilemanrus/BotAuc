package storage

import (
	"context"
)

type Auction struct {
	Name      string
	StartDate string
	EndDate   string
	URL       string
	Status    string
}

type UrlsAlias []string

type Storage interface {
	SaveData(ctx context.Context, p *Auction) error
	IsExists(ctx context.Context, p *Auction) (bool, error)
	ActualizeAucs(ctx context.Context, urls *UrlsAlias) error
	GetFutureAucs(ctx context.Context) (string, error)
}
