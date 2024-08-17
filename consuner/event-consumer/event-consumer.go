package event_consumer

import (
	"context"
	"goTelegram/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(context.Background(), c.batchSize)
		if err != nil {
			log.Printf("can't fetch events: %v", err)

			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(time.Second * 1)

			continue
		}
		if err := c.handleEvents(context.Background(), gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		log.Printf("got event: %v", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't process event: %v", err)

			continue
		}
	}
	return nil
}
