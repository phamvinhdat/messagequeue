package messagequeue

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	unsupportedContentType = errors.New("unsupported content type")
	publisherIsClosed      = errors.New("publisher is closed")
)

type Publisher interface {
	Publish(ctx context.Context, data interface{}, opts ...PublishOption) error
	Close() error
}

type publisher struct {
	sender   Sender
	isClosed bool
}

// Publish public a message, if content-type not assign,
// json-content-type will assign
func (p *publisher) Publish(ctx context.Context, data interface{},
	opts ...PublishOption) error {
	if p.isClosed {
		return publisherIsClosed
	}

	msg, err := buildMsg(data, opts...)
	if err != nil {
		return err
	}

	return p.sender.Send(ctx, msg)
}

func (p *publisher) Close() error {
	if p.isClosed {
		return nil
	}
	p.isClosed = true

	return p.sender.Close()
}

func NewPublisher(sender Sender) Publisher {
	return &publisher{
		sender:   sender,
		isClosed: false,
	}
}

func buildMsg(data interface{}, opts ...PublishOption) (msg Message, err error) {
	opt := GetPublishOption(opts...)

	contentType, ok := opt.Header.Contains(ContentType)
	if !ok {
		contentType = JsonContentType
		opt.Header.Set(ContentType, JsonContentType)
	}

	// convert data
	var bytes []byte
	switch contentType {
	case JsonContentType:
		bytes, err = json.Marshal(data)
		if err != nil {
			return
		}
	case BinaryContentType:
		bytes = data.([]byte)
	default:
		err = unsupportedContentType
		return
	}

	msg.Header = opt.Header
	msg.Data = bytes
	msg.TimeStamp = time.Now()

	return
}
