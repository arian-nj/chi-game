package rooms

import (
	"errors"
	"time"
)

const MaxMessageLen = 1000

type Message struct {
	ID     string
	RoomID string
	UserID string
	Text   string
	SentAt time.Time
}

func (m *Message) Validate() error {
	if len(m.Text) == 0 {
		return errors.New("message text can not be empty")
	}
	if len(m.Text) > MaxMessageLen {
		return errors.New("message to long")
	}
	return nil
}

func (m *Message) CanEdit() bool {
	return time.Since(m.SentAt) < 2*time.Hour
}
