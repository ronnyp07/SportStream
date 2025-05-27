package cccorrelation

import "context"

type Context interface {
	isContextContrain()
}

type DefaultContext struct {
	context.Context
}

func (DefaultContext) isContextContrain() {}
