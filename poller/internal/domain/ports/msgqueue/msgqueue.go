package msgqueue

import (
	"context"

	"github.com/ronnyp07/SportStream/internal/domain/services/msgqueue/msgtype"
)

type MsgQueue interface {
	Connect(ctx context.Context) error
	Close()
	IsConnected() bool
	PublishMessage(ctx context.Context, msgType msgtype.MessageType, data []byte) error
}
