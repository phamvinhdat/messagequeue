package rmq

import (
	"context"
	"time"

	"github.com/go-stomp/stomp/frame"
	"github.com/phamvinhdat/messagequeue"
	"github.com/phamvinhdat/messagequeue/publishoption"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (q *msgQueue) Send(ctx context.Context, msg messagequeue.Message) error {
	err := q.exchangeDeclare("topic")
	if err != nil {
		return err
	}
	header := convertToAMQPHeader(msg.Header)
	contentType := msg.Header.Get(publishoption.ContentType)
	logrus.Debugf("send message to queue, topic: %s, key: %s", q.conf.Topic, q.conf.Key)
	err = q.sendMessage(msg.Data, header, contentType, msg.TimeStamp)

	if err != nil {
		return err
	}

	return nil
}

func (q *msgQueue) sendMessage(data []byte, header amqp.Table, contentType string, timeStamp time.Time) error {
	err := q.channel.Publish(
		q.conf.Topic,
		q.conf.Key,
		false,
		false,
		amqp.Publishing{
			Headers:      header,
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			Timestamp:    timeStamp,
			Body:         data,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *msgQueue) exchangeDeclare(kind string) error {
	err := s.channel.ExchangeDeclare(
		s.conf.Topic, // exchange name
		kind,         // exchange type
		true,         //durable
		false,        // auto delete
		false,        // internal
		false,        // noWait
		nil,          //args
	)
	if err != nil {
		return err
	}
	return nil
}

func convertToAMQPHeader(header frame.Header) amqp.Table {
	table := amqp.Table{}
	for i := 0; i < header.Len(); i++ {
		k, v := header.GetAt(i)
		table[k] = v
	}

	return table
}
