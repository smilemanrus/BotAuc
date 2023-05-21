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
	"time"
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
	StartDate time.Time
	EndDate   time.Time
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
	if err := aucDataByStatus(&aucs, aucReady); err != nil {
		return aucs, e.WrapIfErr(fmt.Sprintf(errLog, aucReady), err)
	}

	if err := aucDataByStatus(&aucs, aucActive); err != nil {
		return aucs, e.WrapIfErr(fmt.Sprintf(errLog, aucActive), err)
	}
	return aucs, nil
}

func aucToEvent(auc Auction, status string) (events.Event, error) {
	res := events.Event{
		Type: events.Auction,
		Text: auc.Name,
	}
	layout := "02.01.2006 15:04"
	frmtdStartDate, err := time.Parse(layout, strings.TrimSpace(auc.StartDate))
	if err != nil {
		return res, e.WrapIfErr("can't parse start date ", err)
	}
	frmtdEndDate, err := time.Parse(layout, strings.TrimSpace(auc.EndDate))
	if err != nil {
		return res, e.WrapIfErr("can't parse end date", err)
	}

	res.Meta = Meta{
		StartDate: frmtdStartDate,
		EndDate:   frmtdEndDate,
		URL:       auc.URL,
		Status:    status,
	}
	return res, err
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

func aucDataByStatus(events *[]events.Event, status string) error {

	c := colly.NewCollector()
	aucs := make([]Auction, 0)

	c.OnHTML("div.panel.panel-default.table-responsive tbody>tr ", func(e *colly.HTMLElement) {

		href, _ := e.DOM.Find("td:nth-child(1)>a").Attr("href")
		href = path.Join(domain, href)

		auc := NewAuc(
			e.DOM.Find("td:nth-child(1)").Text(),
			e.DOM.Find("td:nth-child(3)").Text(),
			e.DOM.Find("td:nth-child(4)").Text(),
			fmt.Sprintf("%s%s", "https://", href))
		aucs = append(aucs, auc)
	})

	url := path.Join(domain, "auction", status)
	if err := c.Visit(fmt.Sprintf("%s%s", "https://", url)); err != nil {
		return e.Wrap("can't parse site", err)
	}

	for _, auc := range aucs {
		aucEvent, err := aucToEvent(auc, status)
		if err != nil {
			return e.Wrap("can't convert auc to event", err)
		}
		*events = append(*events, aucEvent)
	}
	return nil
}

func (p *Processor) actualizeAucs(actualURLS *storage.UrlsAlias) error {
	err := p.storage.ActualizeAucs(context.Background(), actualURLS)
	err = e.WrapIfErr("can't actualise aucs", err)
	return err
}
