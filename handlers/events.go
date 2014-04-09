package handlers

import (
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
	"log"
)

type Event struct {
	Name string
	Body []byte
}

type Events []*Event

func (e Events) Publish() {
	ch := getChannel()
	defer ch.Close()
	for _, event := range e {
		event.PublishOnChannel(ch)
	}
}

func (e *Event) PublishOnChannel(ch *amqp.Channel) {
	if ch == nil {
		ch = getChannel()
		defer ch.Close()
	}

	ch.ExchangeDeclare(
		config.MainExchange, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)

	ch.Publish(
		config.MainExchange, // exchange
		e.Name,              // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        e.Body,
		})
}

func getChannel() *amqp.Channel {
	ch, err := connection.Channel()
	if err != nil {
		dialRMQ()
		ch, err = connection.Channel()
		if err != nil {
			log.Println("rmq", err.Error())
		}
	}

	return ch
}
