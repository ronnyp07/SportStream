package ccnats

import "github.com/nats-io/nats.go/jetstream"

type consumeMessage struct {
	jsMsg jetstream.Msg
}

func (m consumeMessage) Headers() map[string][]string {
	return m.jsMsg.Headers()
}

func (m consumeMessage) Subject() string {
	return m.jsMsg.Subject()
}

func (m consumeMessage) Data() []byte {
	return m.jsMsg.Data()
}

func (m consumeMessage) Ack() error {
	return m.jsMsg.Ack()
}

func (m consumeMessage) Nack() error {
	return m.jsMsg.Nak()
}
