package telegram

import (
	"context"
	"errors"
	"goTelegram/client/telegram"
	err2 "goTelegram/lib/err"
	"goTelegram/storage"
	"log"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got command '%s' from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCmd)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = err2.WrapIfError(err, "can't save page")
	}()

	sendMessage := newMessageSender(chatID, p.tg)

	page := &storage.Page{
		Url:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return sendMessage(msgAlready)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		err = err2.WrapIfError(err, "can't pick random page")
	}()

	sendMessage := newMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMessage(msgNoSaved)
	}

	if err := p.tg.SendMessage(chatID, page.Url); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func newMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	return strings.HasPrefix(text, "https://")
}
