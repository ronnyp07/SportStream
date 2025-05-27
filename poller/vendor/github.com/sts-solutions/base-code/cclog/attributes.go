package cclog

import "go.uber.org/zap/zapcore"

// Attributes add extra key value pairs which will be in the logs.
// They can be used to add context information to the logs like user id, bet id, etc which are not
// parameters of the current event to be logged.
type Attributes struct {
	Fields []zapcore.Field
}

// NewAttributes creates an attribute set for logger.
// It is useful for passing non-volatile parameters to logger like span id, trace id, or
// fields which will be used as Attributes in the log messages. Anything which we pass
// here, will be merged with the fields passed in like logger_info and appear as Attributes
// in Open Telemetry log record.
func NewAttributes(fields ...zapcore.Field) Attributes {
	return Attributes{
		Fields: fields,
	}
}

// Set appends or overwrites a field in the Attributes field list.
func (attrs *Attributes) Set(field zapcore.Field) {
	for i := range attrs.Fields {
		if attrs.Fields[i].Key == field.Key {
			attrs.Fields[i] = field
			return
		}
	}
	attrs.Fields = append(attrs.Fields, field)
}

// Get a reference to the field in the Attributes set or nil if no field with such a key.
func (attrs *Attributes) Get(fieldKey string) *zapcore.Field {
	for i := range attrs.Fields {
		if attrs.Fields[i].Key == fieldKey {
			return &attrs.Fields[i]
		}
	}
	return nil
}

// MarshalLogObject is an internal function implementing object serialization in zap.
func (attrs *Attributes) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for i := 0; i < len(attrs.Fields); i++ {
		attrs.Fields[i].AddTo(enc)
	}
	return nil
}
