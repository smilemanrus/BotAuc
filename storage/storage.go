package storage

import (
	"BotAuc/lib/e"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	SaveData(ctx context.Context, p *Auc) error
	RemoveData(ctx context.Context, p *Auc) error
}

type Auc struct {
	URL       string
	AucName   string
	StartDate string
	EndDate   string
	Id        string
}

func (p Auc) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash (URL)", err)
	}
	if _, err := io.WriteString(h, p.AucName); err != nil {
		return "", e.Wrap("can't calculate hash (UserName)", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
