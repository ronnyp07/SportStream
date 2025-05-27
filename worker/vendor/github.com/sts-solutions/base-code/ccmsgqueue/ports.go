package ccmsgqueue

import "context"

type Connection interface {
	Connect() error
	Close()
	IsConnected() bool
}

type Consumer interface {
	Consume(ctx context.Context) error
	Close(ctx context.Context)
}

type Publisher interface {
	Publish(ctx context.Context, msg PublishMessage) error
}

type ConsumeMessage interface {
	Headers() map[string][]string
	Subject() string
	Data() []byte
	Ack() error
	Nack() error
}

type PublishMessage interface {
	Headers() map[string][]string
	Subject() string
	Data() []byte
}
