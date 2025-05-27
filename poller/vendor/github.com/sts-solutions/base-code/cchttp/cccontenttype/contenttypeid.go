package cccontenttype

type contentTypeID string

func (c contentTypeID) String() string {
	return string(Key)
}

const (
	Key contentTypeID = "Content-Type"
)
