package messagequeue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/phamvinhdat/messagequeue/publishoption"
)

var (
	invalidContentType     = errors.New("invalid content-type")
	unsupportedContentType = errors.New("unsupported content type")
	publisherIsClosed      = errors.New("publisher is closed")
)

type Publisher interface {
	Publish(ctx context.Context, data interface{}, opts ...publishoption.PublishOption) error
	Close() error
}

type publisher struct {
	sender   Sender
	isClosed bool
}

func (p *publisher) Publish(ctx context.Context, data interface{},
	opts ...publishoption.PublishOption) error {
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

func buildMsg(data interface{}, opts ...publishoption.PublishOption) (msg Message, err error) {
	opt := publishoption.GetPublishOption(opts...)

	contentType, ok := opt.Header.Contains(publishoption.ContentType)
	if !ok {
		err = invalidContentType
		return
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
