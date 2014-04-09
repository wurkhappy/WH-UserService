package handlers

import (
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
)

var connection *amqp.Connection

func Setup() {
	dialRMQ()
}

func dialRMQ() {
	var err error
	connection, err = amqp.Dial(config.RMQBroker)
	if err != nil {
		panic(err)
	}
}
