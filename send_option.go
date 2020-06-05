package messagequeue

import "github.com/go-stomp/stomp/frame"

const (
	ContentType = "Content-Type"
)

type PublishOption interface {
	apply(opt *publishOpt)
}

type publishOpt struct {
	Header frame.Header
}

type sendOptionFn func(opt *publishOpt)

func (f sendOptionFn) apply(opt *publishOpt) {
	f(opt)
}

func WithHeader(key, value string) PublishOption {
	return sendOptionFn(func(opt *publishOpt) {
		opt.Header.Set(key, value)
	})
}

func WithHeaderAdding(key, value string) PublishOption {
	return sendOptionFn(func(opt *publishOpt) {
		opt.Header.Add(key, value)
	})
}

func WithContentType(value string) PublishOption {
	return WithHeader(ContentType, value)
}

func GetPublishOption(opts ...PublishOption) publishOpt {
	opt := publishOpt{Header: frame.Header{}}
	for _, optFn := range opts {
		optFn.apply(&opt)
	}

	return opt
}
