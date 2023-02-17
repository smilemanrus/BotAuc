package storage

import (
	"BotAuc/lib/e"
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash (URL)", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash (UserName)", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
