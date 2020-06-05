package msghandler

import (
	"github.com/phamvinhdat/messagequeue"
)

type HookFn func(ctx *QueueContext)

type MsgHandler interface {
	Use(hookFn ...HookFn)
	Subscribe() error
	Close()
}

type msgHandler struct {
	consumer messagequeue.Consumer
	hookFns  []HookFn
	closeCh  chan bool
	isClose  bool
}

func (m *msgHandler) Close() {
	if m.isClose {
		return
	}

	m.isClose = true
	m.closeCh <- true
}

func (m *msgHandler) handle(message messagequeue.Message) {
	ctx := newQueueContext(m.hookFns, message)
	ctx.Next()
}

func (m *msgHandler) Use(hookFn ...HookFn) {
	m.hookFns = append(m.hookFns, hookFn...)
}

func (m *msgHandler) Subscribe() error {
	ch, err := m.consumer.Consume()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				m.handle(msg)
			case <-m.closeCh:
				return
			}
		}
	}()

	return nil
}

func New(consumer messagequeue.Consumer, hookFns ...HookFn) MsgHandler {
	return &msgHandler{
		consumer: consumer,
		hookFns:  hookFns,
		closeCh:  make(chan bool),
	}
}
