package cccorrelation

import "context"

type correlationId string

var CorrelationKey = correlationId("CorrelationId")

func WithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, CorrelationKey, correlationId)
}

func GetCorrelationId(ctx context.Context) (bool, string) {
	val := ctx.Value(CorrelationKey)
	if val == nil {
		return false, ""
	}

	return true, val.(string)
}
