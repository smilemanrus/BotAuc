package storage

import (
	"context"
)

type Auction struct {
	Name      string
	StartDate string
	EndDate   string
	URL       string
}

type Storage interface {
	SaveData(ctx context.Context, p *Auction) error
	RemoveData(ctx context.Context, p *Auction) error
	IsExists(ctx context.Context, p *Auction) (bool, error)
}
