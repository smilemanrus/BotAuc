package auctions

import (
	"BotAuc/events"
	"github.com/gocolly/colly/v2"
	"strings"
)

const (
	url = "https://auctions.partner.ru/auction/ready"
)

type Processor struct {
	offset         int
	ProcessMessage interface{}
}

type Auction struct {
	Name      string
	StartDate string
	EndDate   string
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
}

func New() *Processor {
	return &Processor{}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {

	aucs := make([]events.Event, 0)
	c := colly.NewCollector()

	c.OnHTML("div.panel.panel-default.table-responsive tbody>tr ", func(e *colly.HTMLElement) {
		auc := NewAuc(
			e.DOM.Find("td:nth-child(1)").Text(),
			e.DOM.Find("td:nth-child(3)").Text(),
			e.DOM.Find("td:nth-child(4)").Text())
		aucEvent := event(auc)
		aucs = append(aucs, aucEvent)
	})

	err := c.Visit(url)

	return aucs, err
}

func event(auc Auction) events.Event {
	res := events.Event{
		Type: events.Auction,
		Text: auc.Name,
	}
	res.Meta = Meta{
		StartDate: auc.StartDate,
		EndDate:   auc.EndDate,
	}
	return res
}
func (p *Processor) Process(event events.Event) error {

	return nil
}
