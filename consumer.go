package messagequeue

type Consumer interface {
	Consume() (<-chan Message, error)
	Close() error
}
