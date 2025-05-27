package ccnats

type publishMessage struct {
	headers map[string][]string
	data    []byte
	subject string
}

func (m *publishMessage) Headers() map[string][]string {
	return m.headers
}

func (m *publishMessage) Subject() string {
	return m.subject
}

func (m *publishMessage) Data() []byte {
	return m.data
}
