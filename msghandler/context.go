package msghandler

import (
	"context"

	"github.com/phamvinhdat/messagequeue"
)

type QueueContext struct {
	context.Context

	message    messagequeue.Message
	err        error
	hookFns    []HookFn
	index      int
	abortIndex int
	result     interface{}
}

func newQueueContext(hookFns []HookFn, message messagequeue.Message) *QueueContext {
	return &QueueContext{
		Context:    context.Background(),
		message:    message,
		hookFns:    hookFns,
		abortIndex: -1,
		index:      -1,
	}
}

func (c *QueueContext) Next() {
	if c.isAbort() {
		return
	}

	lenHookFn := len(c.hookFns)
	c.index++
	if lenHookFn <= c.index {
		return
	}

	hookFn := c.hookFns[c.index]
	hookFn(c)
}

func (c *QueueContext) Error() error {
	return c.err
}

func (c *QueueContext) Message() messagequeue.Message {
	return c.message
}

func (c *QueueContext) Result() interface{} {
	return c.result
}

func (c *QueueContext) AbortWithError(err error) {
	if err == nil {
		panic("err is nil")
	}

	c.Abort()
	c.err = err
}

func (c *QueueContext) AbortWithResult(err error, result interface{}) {
	c.AbortWithError(err)
	c.result = result
}

func (c *QueueContext) Abort() {
	c.abortIndex = c.index
}

func (c *QueueContext) WithValue(key, value interface{}) {
	if c.Context == nil {
		c.Context = context.Background()
	}

	c.Context = context.WithValue(c.Context, key, value)
}

func (c *QueueContext) isAbort() bool {
	return c.abortIndex > -1
}
