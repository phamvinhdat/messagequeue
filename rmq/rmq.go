package rmq

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	URL   string
	Topic string
	Key   string
}

type msgQueue struct {
	conf      *RabbitMQConfig
	conn      *amqp.Connection
	channel   *amqp.Channel
	errorChan chan *amqp.Error
	isClosed  bool
}

func New(conf *RabbitMQConfig) (*msgQueue, error) {
	mQueue := &msgQueue{
		conf: conf,
	}

	err := mQueue.connect()
	if err != nil {
		return nil, err
	}
	go mQueue.reconnector()
	return mQueue, nil
}

func (q *msgQueue) Close() error {
	if q.isClosed {
		return nil
	}
	logrus.Info("close connection")
	q.isClosed = true
	if q.channel != nil {
		_ = q.channel.Close()
		q.channel = nil
	}

	if q.conn != nil {
		if err := q.conn.Close(); err != nil {
			return err
		}
		q.conn = nil
	}
	return nil
}

func (q *msgQueue) connect() error {
	err := q.createConnection()
	if err != nil {
		return err
	}

	err = q.createChannel()
	if err != nil {
		_ = q.conn.Close()
		return err
	}

	q.errorChan = make(chan *amqp.Error)
	q.conn.NotifyClose(q.errorChan)
	logrus.Info("connection established")
	return nil
}

func (q *msgQueue) reconnect() {
	for {
		err := q.connect()
		if err == nil {
			logrus.Info("connection established")
			return
		}

		logrus.Error("failed to reconnect, retrying in 1 second, error: ", err)
		time.Sleep(time.Second)
	}
}

func (q *msgQueue) reconnector() {
	for {
		err := <-q.errorChan
		if q.isClosed {
			return
		}

		logrus.Error("trying reconnect, error: ", err)
		q.reconnect()
	}
}

func (q *msgQueue) createConnection() error {
	conn, err := amqp.Dial(q.conf.URL)
	if err != nil {
		return err
	}

	q.conn = conn
	return nil
}

func (q *msgQueue) createChannel() error {
	channel, err := q.conn.Channel()
	if err != nil {
		return err
	}

	q.channel = channel
	return nil
}
