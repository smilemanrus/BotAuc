package auctions

import (
	"BotAuc/events"
	"BotAuc/lib/e"
	"BotAuc/storage"
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"path"
	"strings"
)

const (
	domain    = "auctions.partner.ru"
	aucReady  = "ready"
	aucActive = "active"
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

func NewAuc(name, startDate, endDate, url string) Auction {
	return Auction{
		Name:      strings.TrimSpace(name),
		StartDate: strings.TrimSpace(startDate),
		EndDate:   strings.TrimSpace(endDate),
		URL:       strings.TrimSpace(url),
	}
}

type Meta struct {
	StartDate string
	EndDate   string
	URL       string
	Status    string
}

func New(storage storage.Storage) *Processor {
	return &Processor{
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	aucs := make([]events.Event, 0)
	errLog := "can't parse auc, status %s"
	err := aucDataByStatus(&aucs, aucReady)
	err = e.WrapIfErr(fmt.Sprintf(errLog, aucReady), err)
	err = aucDataByStatus(&aucs, aucActive)
	err = e.WrapIfErr(fmt.Sprintf(errLog, aucActive), err)

	return aucs, err
}

func aucToEvent(auc Auction, status string) events.Event {
	res := events.Event{
		Type: events.Auction,
		Text: auc.Name,
	}
	res.Meta = Meta{
		StartDate: auc.StartDate,
		EndDate:   auc.EndDate,
		URL:       auc.URL,
		Status:    status,
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
		Status:    meta.Status,
	}
	return auc, nil
}
func (p *Processor) Process(events []events.Event) error {
	var err error
	var actURL string
	actualURLS := make(storage.UrlsAlias, 0)
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)
		actURL, err = p.SaveEvent(event)
		if err != nil {
			continue
		}
		actualURLS = append(actualURLS, actURL)
	}

	if err = p.actualizeAucs(&actualURLS); err != nil {
		err = e.Wrap("can't process auc actualising", err)
	}

	return err
}

func (p *Processor) SaveEvent(event events.Event) (string, error) {
	errMsg := "can't process event"
	auc, err := eventToAuc(event)
	err = e.WrapIfErr(errMsg, err)

	isExist, err := p.storage.IsExists(context.Background(), &auc)
	err = e.WrapIfErr(errMsg, err)

	if !isExist {
		err = p.storage.SaveData(context.Background(), &auc)
		err = e.WrapIfErr(errMsg, err)
	}
	return auc.URL, err
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", events.ErrUnknownMetaType())
	}
	return res, nil
}

func aucDataByStatus(aucs *[]events.Event, status string) error {

	c := colly.NewCollector()
	c.OnHTML("div.panel.panel-default.table-responsive tbody>tr ", func(e *colly.HTMLElement) {

		href, _ := e.DOM.Find("td:nth-child(1)>a").Attr("href")
		href = path.Join(domain, href)

		auc := NewAuc(
			e.DOM.Find("td:nth-child(1)").Text(),
			e.DOM.Find("td:nth-child(3)").Text(),
			e.DOM.Find("td:nth-child(4)").Text(),
			fmt.Sprintf("%s%s", "https://", href))
		aucEvent := aucToEvent(auc, status)
		*aucs = append(*aucs, aucEvent)
	})

	url := path.Join(domain, "auction", status)
	err := c.Visit(fmt.Sprintf("%s%s", "https://", url))
	return err
}

func (p *Processor) actualizeAucs(actualURLS *storage.UrlsAlias) error {
	err := p.storage.ActualizeAucs(context.Background(), actualURLS)
	err = e.WrapIfErr("can't actualise aucs", err)
	return err
}
