package ccotelnats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer oteltrace.Tracer
var conn *nats.Conn

func SetTracer(nc *nats.Conn, tr oteltrace.Tracer) {
	tracer = tr
	conn = nc
}

func StartPublishSpan(ctx context.Context, msg *nats.Msg) (context.Context, oteltrace.Span) {
	carrier := propagation.HeaderCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	headers := nats.Header(carrier)

	if msg.Header == nil {
		msg.Header = make(nats.Header)
	}

	for k, v := range headers {
		msg.Header[k] = v
	}

	return tracer.Start(ctx,
		fmt.Sprintf("PUBLISH: %v", msg.Subject),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(getCommonTraceAttributes()...),
		oteltrace.WithAttributes([]attribute.KeyValue{
			semconv.MessagingOperationKey.String("PUBLISH"),
			semconv.MessagingDestinationKindKey.String(msg.Subject),
			semconv.MessagingOperationKey.String("SUBJECT"),
			semconv.MessagingMessagePayloadCompressedSizeBytesKey.Int(len(msg.Data)),
		}...),
	)
}

func StartSubscriberSpan(ctx context.Context, msg *nats.Msg, group string) (context.Context, oteltrace.Span) {
	return tracer.Start(otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(msg.Header)),
		fmt.Sprintf("SUBSCRIBE %v", msg.Subject),
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(getCommonTraceAttributes()...),
		oteltrace.WithAttributes([]attribute.KeyValue{
			semconv.MessagingOperationKey.String("SUBSCRIBE"),
			semconv.MessagingMessagePayloadCompressedSizeBytesKey.Int(len(msg.Data)),
			attribute.String("messaging.source", msg.Subject),
			attribute.String("messaging.source_group", group),
		}...),
	)
}

func getCommonTraceAttributes() []attribute.KeyValue {
	if conn == nil {
		return []attribute.KeyValue{}
	}

	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String("NATS"),
		semconv.MessagingSystemKey.String(conn.ConnectedAddr()),
	}
}

func StartSubscriberSpanJetStream(msg jetstream.Msg, group string) (context.Context, oteltrace.Span) {
	return tracer.Start(otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(msg.Headers())),
		fmt.Sprintf("SUBSCRIBE %v", msg.Subject),
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(getCommonTraceAttributes()...),
		oteltrace.WithAttributes([]attribute.KeyValue{
			semconv.MessagingOperationKey.String("SUBSCRIBE"),
			semconv.MessagingMessagePayloadCompressedSizeBytesKey.Int(len(msg.Data())),
			attribute.String("messaging.source", msg.Subject()),
			attribute.String("messaging.source_group", group),
		}...),
	)
}
