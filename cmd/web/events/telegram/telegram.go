package telegram

import (
	"errors"

	"taskbot/cmd/web/clients/telegram"
	"taskbot/cmd/web/storage"
	"taskbot/internal/e"
)

type Processor struct {
	tg *telegram.Client
	offset int
	storage storage.Storage
}

type Meta struct {
	ChatID	 int,
	Username string,
}

var ( 
	ErrUnknownEventType = errors.New("Unknown event type.")
	ErrUnknownMetaType = errors.New("Unknown meta type.")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor {
		tg: client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.WrapIfError("can't get events" , err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u)):w
	}

	p.offset = updates[len(updates) - 1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.WrapIfErr("Can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.WrapIfErr("Can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.WrapIfErr("Can't process message", err)
	}

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meat)
	if !ok {
		return Meta{}, e.WrapIfErr("Can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(update telegram.Update) events.Event {
	updateType = fetchType(update)
	res := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		res.Meta = Meta {
			ChatID: update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
 
}

func fetchText(update telegram.Uodate) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}

func fetchType(update telegram.Uodate) string {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message.
}
