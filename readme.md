# Abstraction message queue for go

## support

- [x] [RabbitMQ](https://www.rabbitmq.com/)
- [ ] [Kafka](https://kafka.apache.org/) (coming soon) 
- [ ] [ActiveMQ](https://activemq.apache.org/) (coming soon)

## Feature
- [X] Auto connect when losing connection
- [X] Abstraction publisher and consumer:
    ```go
    type Publisher interface {
          Publish(ctx context.Context, data interface{}, opts ...publishoption.PublishOption) error	 
          Close() error
    }
  
    type Consumer interface {
    	Consume() (<-chan Message, error)
    	Close() error
    }
    ```
- [X] Easy to use

## Example

### Publisher

```shell script
go run main.go "yourKey"
```
```go
package main

import (
	"context"
	"os"

	"github.com/phamvinhdat/messagequeue"
	"github.com/phamvinhdat/messagequeue/publishoption"
	"github.com/phamvinhdat/messagequeue/rmq"
)

var conf = rmq.RabbitMQConfig{
	URL:      "amqp://guest:guest@localhost:5672/",
	Exchange: "logs_topic",
	Key:      severityFrom(os.Args),
}

func main() {
	messagequeue.WithDebug()
	rmq, err := rmq.New(&conf)
	if err != nil {
		panic(err)
	}
	defer rmq.Close()

	publisher := messagequeue.NewPublisher(rmq)
	defer publisher.Close()

	err = publisher.Publish(context.Background(), conf, publishoption.WithContentType(messagequeue.JsonContentType))
	if err != nil {
		panic(err)
	}
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
```

### Consumer

```go
package main

import (
	"log"
	"os"
	
	"github.com/phamvinhdat/messagequeue"
	"github.com/phamvinhdat/messagequeue/rmq"
)

func main() {
	var conf = rmq.RabbitMQConfig{
		URL:      "amqp://guest:guest@localhost:5672/",
		Exchange: "logs_topic",
		Key:      severityFrom(os.Args),
	}
	messagequeue.WithDebug()
	rmq, err := rmq.New(&conf)
	if err != nil {
		panic(err)
	}
	defer rmq.Close()

	msgChan, err := rmq.Consume()
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgChan {
			log.Println(d)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}

```
