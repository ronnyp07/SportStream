package msgtype

type MessageType int32

const (
	NotSet MessageType = iota
	ArticlesUpdated
)

var namesMessageType = map[MessageType]string{
	NotSet:          `not-set`,
	ArticlesUpdated: `article-updated`,
}

func FromName(name string) MessageType {
	for id, n := range namesMessageType {
		if n == name {
			return id
		}
	}

	return NotSet
}

func Get(id int32) MessageType {
	v := MessageType(id)

	if _, ok := namesMessageType[v]; ok {
		return v
	}

	return NotSet
}

func (t MessageType) ID() int32 {
	return int32(t)
}
