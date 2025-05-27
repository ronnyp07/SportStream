package ccnats

import (
	"context"

	"github.com/sts-solutions/base-code/ccotel/ccotelnats"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
)

// Tracer provides tracing capabilities for NATS operations
type Tracer struct {
	tracer      trace.Tracer
	shouldTrace bool
}

// NewTracer creates a new Tracer instance
func NewTracer(tracer trace.Tracer) Tracer {
	return Tracer{
		tracer:      tracer,
		shouldTrace: tracer != nil,
	}
}

// SetConnectionTracer configures tracing for a NATS connection
func (t Tracer) setConnectionTracer(conn *nats.Conn) {
	if t.shouldTrace {
		ccotelnats.SetTracer(conn, t.tracer)
	}
}

// StartConsumerSpan starts a new span for message consumption
func (t Tracer) startConsumerSpan(jsMsg jetstream.Msg, name string) (context.Context, func(...trace.SpanEndOption)) {
	if t.shouldTrace {
		spanCtx, span := ccotelnats.StartSubscriberSpanJetStream(jsMsg, name)
		return spanCtx, span.End
	}
	return context.Background(), func(...trace.SpanEndOption) {}
}

// StartPublisherSpan starts a new span for message publishing
func (t Tracer) startPublisherSpan(msg *nats.Msg, ctx context.Context) (context.Context, func(...trace.SpanEndOption)) {
	if t.shouldTrace {
		spanCtx, span := ccotelnats.StartPublishSpan(ctx, msg)
		return spanCtx, span.End
	}
	return context.Background(), func(...trace.SpanEndOption) {}
}
