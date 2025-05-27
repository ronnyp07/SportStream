package msgqueue

const (
	MessageHeaders = "msg_headers"
	MessageSubject = "msg_subject"
	MessageData    = "msg_data"
)

type QueueMessage struct {
	Header  map[string][]string
	Subject string
	Data    []byte
}

func (q QueueMessage) AsJson() map[string]any {
	return map[string]any{
		MessageHeaders: q.Header,
		MessageSubject: q.Subject,
		MessageData:    q.Data,
	}
}
