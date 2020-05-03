package messagequeue

import "context"

type Sender interface {
	Send(ctx context.Context, msg Message) error
	Close() error
}
