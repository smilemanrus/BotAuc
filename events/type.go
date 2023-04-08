package events

import "errors"

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	Auction
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

func ErrUnknownMetaType() error {
	return errors.New("unknown meta type")
}
