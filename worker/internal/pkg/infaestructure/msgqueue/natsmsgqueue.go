package msgqueue

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/ronnyp07/SportStream/worker/internal/metrics/promnats"
	cnats "github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/nats"
	"github.com/sts-solutions/base-code/ccotel/ccotelnats"

	"emperror.dev/errors"
)

type MsgQueueService interface {
	PublishMessage(ctx context.Context, subject string, data []byte) (
		msg QueueMessage, err error)
}

type natsService struct {
	natsJSCtx nats.JetStreamContext
}

func NewNatsService(jsCtxt nats.JetStreamContext) MsgQueueService {
	return &natsService{
		natsJSCtx: jsCtxt,
	}
}

func (n natsService) PublishMessage(ctx context.Context, subject string, data []byte) (
	msg QueueMessage, err error) {

	natsMsg := &nats.Msg{
		Subject: subject,
		Data:    data,
	}

	msg.Subject = subject
	msg.Data = data
	msg.Header = make(map[string][]string)
	for k, v := range natsMsg.Header {
		msg.Header[k] = v
	}

	ctx, span := ccotelnats.StartPublishSpan(ctx, natsMsg)
	cnats.SetMsgCorrelationId(ctx, natsMsg)
	defer span.End()

	_, err = promnats.PublishMsg(n.natsJSCtx, natsMsg)
	if err != nil {
		return msg, errors.Wrap(err, "publishing message in nats queue")
	}

	return msg, err
}
