package app

import (
	"github.com/ronnyp07/SportStream/internal/pkg/infaestructure/msgqueue"
)

type Services struct {
	MsgQueueService msgqueue.MsgQueueService
}
