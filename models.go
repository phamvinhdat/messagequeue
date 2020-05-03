package messagequeue

import (
	"time"

	"github.com/go-stomp/stomp/frame"
)

const (
	JsonContentType   = "application/json"
	BinaryContentType = "application/octet-stream"
)

type Message struct {
	Header    frame.Header
	Data      []byte
	TimeStamp time.Time
}
