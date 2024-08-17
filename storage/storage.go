package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	err2 "goTelegram/lib/err"
	"io"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	Url      string
	UserName string
}

var ErrNoSavedPages = errors.New("no saved pages")

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.Url); err != nil {
		return "", err2.Wrap(err, "can't calculate hash")
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", err2.Wrap(err, "can't calculate hash")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
