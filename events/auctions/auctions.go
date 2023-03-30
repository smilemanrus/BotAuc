package auctions

import (
	"BotAuc/events"
	"fmt"
	"github.com/gocolly/colly/v2"
)

const (
	url = "file:///C:/Users/smile/Downloads/%D0%90%D1%83%D0%BA%D1%86%D0%B8%D0%BE%D0%BD.html"
)

type Processor struct {
	offset         int
	ProcessMessage interface{}
}

func New() *Processor {
	return &Processor{}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {

	aucs := make([]events.Event, 0)
	c := colly.NewCollector()

	c.OnHTML("tbody tr td ", func(e *colly.HTMLElement) {
		fmt.Println(e)
	})

	err := c.Visit(url)

	return aucs, err
}

func (p *Processor) Process(event events.Event) error {

	return nil
}
