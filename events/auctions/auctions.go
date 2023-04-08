package auctions

import (
	"BotAuc/events"
	"BotAuc/lib/e"
	"BotAuc/storage"
	"context"
	"github.com/gocolly/colly/v2"
	"strings"
)

const (
	url = "https://auctions.partner.ru/auction/ready"
)

type Processor struct {
	storage storage.Storage
}

type Auction struct {
	Name      string
	StartDate string
	EndDate   string
	URL       string
}

func NewAuc(name, startDate, endDate string) Auction {
	return Auction{
		Name:      strings.TrimSpace(name),
		StartDate: strings.TrimSpace(startDate),
		EndDate:   strings.TrimSpace(endDate),
	}
}

type Meta struct {
	StartDate string
	EndDate   string
	URL       string
}

func New(storage storage.Storage) *Processor {
	return &Processor{
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {

	aucs := make([]events.Event, 0)
	c := colly.NewCollector()

	c.OnHTML("div.panel.panel-default.table-responsive tbody>tr ", func(e *colly.HTMLElement) {
		auc := NewAuc(
			e.DOM.Find("td:nth-child(1)").Text(),
			e.DOM.Find("td:nth-child(3)").Text(),
			e.DOM.Find("td:nth-child(4)").Text())
		aucEvent := aucToEvent(auc)
		aucs = append(aucs, aucEvent)
	})

	err := c.Visit(url)

	return aucs, err
}

func aucToEvent(auc Auction) events.Event {
	res := events.Event{
		Type: events.Auction,
		Text: auc.Name,
	}
	res.Meta = Meta{
		StartDate: auc.StartDate,
		EndDate:   auc.EndDate,
		URL:       auc.URL,
	}
	return res
}

func eventToAuc(event events.Event) (storage.Auction, error) {
	meta, err := meta(event)
	if err != nil {
		return storage.Auction{}, e.Wrap("can't convert event to auc", err)
	}
	auc := storage.Auction{
		Name:      event.Text,
		StartDate: meta.StartDate,
		EndDate:   meta.EndDate,
		URL:       meta.URL,
	}
	return auc, nil
}
func (p *Processor) Process(event events.Event) error {
	errMsg := "can't process event"
	auc, err := eventToAuc(event)
	err = e.WrapIfErr(errMsg, err)

	isExist, err := p.storage.IsExists(context.Background(), &auc)
	err = e.WrapIfErr(errMsg, err)

	if !isExist {
		err = p.storage.SaveData(context.Background(), &auc)
		err = e.WrapIfErr(errMsg, err)
	}
	return err
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", events.ErrUnknownMetaType())
	}
	return res, nil
}
