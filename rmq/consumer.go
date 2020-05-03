package rmq

import (
	"time"

	"github.com/go-stomp/stomp/frame"
	"github.com/phamvinhdat/messagequeue"
	"github.com/phamvinhdat/messagequeue/publishoption"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/streadway/amqp"
)

func (q *msgQueue) Consume() (<-chan messagequeue.Message, error) {
	delivery, err := q.channelConsume()
	if err != nil {
		return nil, err
	}

	msgChan := make(chan messagequeue.Message)
	go q.loopConsumeMsg(delivery, msgChan)
	return msgChan, nil
}

func (q *msgQueue) loopConsumeMsg(delivery <-chan amqp.Delivery,
	msgChan chan<- messagequeue.Message) {
	var err error
	for !q.isClosed {
		select {
		case chanErr := <-q.errorChan:
			logrus.Error("consume receive error: ", chanErr)
			logrus.Info("reconsume")
			for {
				delivery, err = q.channelConsume()
				if err == nil {
					logrus.Info("reconsume success")
					break
				}

				logrus.Info("failed to reconsume, trying in 2s ...")
				time.Sleep(time.Second * 2)
			}

		case amqpMsg := <-delivery:
			msg := convertAMQPMsgToMessage(amqpMsg)
			msgChan <- msg
		}
	}
}

func (q *msgQueue) prepareConsume() error {
	_, err := q.queueDeclare()
	if err != nil {
		return err
	}

	err = q.queueBind()
	if err != nil {
		return err
	}

	return nil
}

func (q *msgQueue) channelConsume() (<-chan amqp.Delivery, error) {
	err := q.prepareConsume()
	if err != nil {
		return nil, err
	}

	return q.channel.Consume(
		q.conf.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (q *msgQueue) queueDeclare() (amqp.Queue, error) {
	return q.channel.QueueDeclare(
		q.conf.Queue,
		false,
		false,
		true,
		false,
		nil,
	)
}

func (q *msgQueue) queueBind() error {
	return q.channel.QueueBind(
		q.conf.Queue,    // name
		q.conf.Key,      // key
		q.conf.Exchange, // exchange
		false,           //noWait
		nil,             // args
	)
}

func convertToFrameHeader(table amqp.Table) frame.Header {
	header := frame.Header{}
	for k, v := range table {
		value, err := cast.ToStringE(v)
		if err != nil {
			logrus.Error("failed to casting header, err: ", err)
		}
		header.Set(k, value)
	}

	return header
}

func convertAMQPMsgToMessage(amqpMsg amqp.Delivery) messagequeue.Message {
	header := convertToFrameHeader(amqpMsg.Headers)
	header.Set(publishoption.ContentType, amqpMsg.ContentType)

	return messagequeue.Message{
		Header:    header,
		Data:      amqpMsg.Body,
		TimeStamp: amqpMsg.Timestamp,
	}
}
