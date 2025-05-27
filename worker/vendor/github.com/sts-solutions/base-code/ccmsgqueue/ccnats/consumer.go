package ccnats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sts-solutions/base-code/ccmsgqueue"
)

const (
	GetMessageMetadataError = "GetMessageMetadataError"
	ConsumeMessageError     = "ConsumeMessageError"
)

// consumer represents a consumer for NATS messages.
type consumer struct {
	namespace  string
	stream     string
	subject    string
	name       string
	conn       *connection
	ctx        jetstream.ConsumeContext
	msgHandler func(ctx context.Context, msg ccmsgqueue.ConsumeMessage)
	metrics    ccmsgqueue.ConsumerMetrics
}

// Consume starts consuming messages from the specified stream and subject.
func (c *consumer) Consume(ctx context.Context) error {
	jsConsumer, err := c.conn.js.Consumer(ctx, c.stream, c.name)
	if err != nil {
		return err
	}

	c.ctx, err = jsConsumer.Consume(func(jsMsg jetstream.Msg) {
		spanCtx, spanEnd := c.conn.tracer.startConsumerSpan(jsMsg, c.name)
		defer spanEnd()

		msg := &consumeMessage{
			jsMsg: jsMsg,
		}

		meta, mErr := jsMsg.Metadata()
		if mErr == nil {
			c.metrics.MessageLag(time.Since(meta.Timestamp).Seconds(), msg.Subject(), c.name)
		} else {
			c.metrics.ErrorInc(GetMessageMetadataError)
		}

		c.msgHandler(spanCtx, msg)
	})

	if err != nil {
		c.metrics.ErrorInc(ConsumeMessageError)
		return err
	}
	return nil
}

// Close stops the consumer and logs the closure.
func (c *consumer) Close(ctx context.Context) {
	if c.ctx != nil {
		c.ctx.Stop()
	}
	c.conn.logger.LogInfo(ctx, fmt.Sprintf("NATS consumer %s closed successfully", c.name))
}
