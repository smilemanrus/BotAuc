package auctions

import (
	"BotAuc/events"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

type Processor struct {
	offset         int
	ProcessMessage interface{}
}

func New() *Processor {
	return &Processor{}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://auctions.partner.ru/auction/ready"},
		ParseFunc: parseMovies,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()

	aucs := make([]events.Event, 0)
	return aucs, nil
}

func parseMovies(g *geziyor.Geziyor, r *client.Response) {

	r.HTMLDoc.Find("table.table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		resultaa := s.Find("td").Map(func(_ int, s *goquery.Selection) string {
			return s.Text()
		})
		fmt.Print(resultaa)
	})

	//fmt.Print(mains)
	//.Each(func(i int, s *goquery.Selection) {
	//	var sessions = strings.Split(s.Find(".shedule_session_time").Text(), " \n ")
	//	sessions = sessions[:len(sessions)-1]
	//
	//	for i := 0; i < len(sessions); i++ {
	//		sessions[i] = strings.Trim(sessions[i], "\n ")
	//	}
	//	if href, ok := s.Find("a.gtm-ec-list-item-movie").Attr("href"); ok {
	//		g.Get(r.JoinURL(href), func(_g *geziyor.Geziyor, _r *client.Response) {
	//			description = _r.HTMLDoc.Find("span.announce p.movie_card_description_inform").Text()
	//			description = strings.ReplaceAll(description, "\t", "")
	//			description = strings.ReplaceAll(description, "\n", "")
	//			description = strings.TrimSpace(description)
	//}
	//		g.Exports <- map[string]interface{}{
	//			"title":        strings.TrimSpace(s.Find("span.movie_card_header.title").Text()),
	//			"subtitle":    strings.TrimSpace(s.Find("span.sub_title.shedule_movie_text").Text()),
	//			"sessions":    sessions,
	//			"description": description,
	//		}
}
func (p *Processor) Process(event events.Event) error {

	return nil
}
