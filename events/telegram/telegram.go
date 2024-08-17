package telegram

import (
	"context"
	"errors"
	"goTelegram/client/telegram"
	"goTelegram/events"
	err2 "goTelegram/lib/err"
	"goTelegram/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func New(tg *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      tg,
		storage: storage,
	}
}

func (p *Processor) Process(_ context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return err2.Wrap(ErrUnknownEventType, "unknown event type")
	}
}

func (p *Processor) Fetch(_ context.Context, limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, err
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(update))

	for _, u := range update {
		res = append(res, event(u))
	}

	p.offset = update[len(update)-1].UpdateID + 1

	return res, nil

}

func event(update telegram.Update) events.Event {
	updType := fetchType(update)
	res := events.Event{
		Type: updType,
		Text: fetchText(update),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message

}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err2.Wrap(err, "can't get meta")
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return err2.Wrap(err, "can't process command")
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, err2.Wrap(storage.ErrNoSavedPages, "unknown meta type")
	}

	return res, nil
}
