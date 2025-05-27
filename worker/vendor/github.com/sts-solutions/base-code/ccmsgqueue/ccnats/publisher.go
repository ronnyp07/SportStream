package ccnats

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/nats-io/nats.go"
	"github.com/sts-solutions/base-code/ccmsgqueue"
)

const (
	PublishMessageError = "PublishMessageError"
)

// publisher is responsible for publishing messages to a JetStream stream.
type publisher struct {
	conn    *connection
	metrics ccmsgqueue.PublisherMetrics
}

// Publish publishes a message to the JetStream stream.
func (p *publisher) Publish(ctx context.Context, msg ccmsgqueue.PublishMessage) error {
	natsMsg := &nats.Msg{
		Subject: msg.Subject(),
		Data:    msg.Data(),
		Header:  msg.Headers(),
	}

	spanCtx, spanEnd := p.conn.tracer.startPublisherSpan(natsMsg, ctx)
	defer spanEnd()

	start := time.Now()

	ack, err := p.conn.js.PublishMsg(spanCtx, natsMsg)
	if err != nil {
		p.metrics.ErrorInc(PublishMessageError)
		return errors.Wrapf(err, "Error publishing message to queue: %v", msg)
	}

	p.metrics.MessageCounterInc(natsMsg.Subject)
	p.metrics.MessageLag(time.Since(start).Seconds(), natsMsg.Subject)

	p.conn.logger.LogInfo(spanCtx, fmt.Sprintf("Message published to stream %s: %v", ack.Stream, msg))

	return nil
}
