package models

import (
	// "database/sql"
	// _ "github.com/bmizerany/pq"
	"github.com/streadway/amqp"
	// "log"
)

var PaymentInfoService string = "http://localhost:3120"
var connection *amqp.Connection
var emailExchange string = "email"
var emailQueue string = "email"

func init() {
	setup()
}

func setup() {
	var err error
	uri := "amqp://guest:guest@localhost:5672/"
	connection, err = amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
}
