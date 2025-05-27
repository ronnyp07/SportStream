package cccorrelation

import (
	"context"
)

type correlationID string

func (c correlationID) String() string {
	return string(c)
}

const (
	//Key is the name used in the headers and logs for CorrelationIj
	Key    correlationID = "X-Correlation-ID"
	LogKey string        = "correlationID"
)

// GetCorrelationID returns the correlationID from the context
func GetCorrelationID(ctx context.Context) string {
	val := ctx.Value(Key)

	if val != nil {
		return val.(string)
	}

	val = ctx.Value(Key.String())
	if val == nil {
		return ""
	}

	return val.(string)

}

// GetCorrelationIDByCustomkey gets the correlation id from the context, searching by given key
func GetCorrelationIDByCustomKey(ctx context.Context, key string) string {
	val := ctx.Value(key)
	if val == nil {
		return ""
	}

	return val.(string)
}
