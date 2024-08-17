package events

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, limit int) ([]Event, error)
}

type Processor interface {
	Process(ctx context.Context, event Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
