package cccontenttype

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type ContentType interface {
	Name() string
	MarshalFunc() MarshalFunc
	UnmarshalFunc() UnmarshalFunc
}

// EncodingContentType represents the type of content encoding
type EncodingContentType int

const (
	NotSet EncodingContentType = iota
	ApplicationJSON
	ApplicationXML
	FormURLEncoded
	TextXML
	TextPlain
)

// marshalFuncFunc encodes any value to bytes
type MarshalFunc func(any) ([]byte, error)

// DecoderFunc decodes bytes to a value
type UnmarshalFunc func([]byte, any) error

type contentType struct {
	Name          string
	marshalFunc   MarshalFunc
	unmarshalFunc UnmarshalFunc
}

var contentTypes = map[EncodingContentType]contentType{
	NotSet: {
		Name:          "",
		marshalFunc:   nil,
		unmarshalFunc: nil,
	},
	ApplicationJSON: {
		Name: "application/json",
		marshalFunc: func(v any) ([]byte, error) {
			return json.Marshal(v)
		},
		unmarshalFunc: func(data []byte, v any) error {
			return json.Unmarshal(data, v)
		},
	},
	ApplicationXML: {
		Name: "application/xml",
		marshalFunc: func(v any) ([]byte, error) {
			return xml.Marshal(v)
		},
		unmarshalFunc: func(data []byte, v any) error {
			return xml.Unmarshal(data, v)
		},
	},
	TextXML: {
		Name:          "text/xml",
		marshalFunc:   nil,
		unmarshalFunc: nil,
	},
	TextPlain: {
		Name: "text/plain",
		marshalFunc: func(v any) ([]byte, error) {
			if s, ok := v.(string); ok {
				return []byte(s), nil
			}
			return nil, fmt.Errorf("expected string, got %T", v)
		},
		unmarshalFunc: func(data []byte, v any) error {
			if s, ok := v.(*string); ok {
				*s = string(data)
				return nil
			}
			return fmt.Errorf("expected *string, got %T", v)
		},
	},
}

// Name returns the name of the ContentType
func (t EncodingContentType) Name() string {
	ct, ok := contentTypes[t]
	if !ok {
		return contentTypes[NotSet].Name
	}
	return ct.Name
}

// marshalFunc returns the marshalFunc function for the ContentType
func (t EncodingContentType) MarshalFunc() MarshalFunc {
	ct, ok := contentTypes[t]
	if !ok {
		return contentTypes[NotSet].marshalFunc
	}
	return ct.marshalFunc
}

// Decoder returns the decoder function for the ContentType
func (t EncodingContentType) UnmarshalFunc() UnmarshalFunc {
	ct, ok := contentTypes[t]
	if !ok {
		return contentTypes[NotSet].unmarshalFunc
	}
	return ct.unmarshalFunc
}
