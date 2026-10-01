package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)





type RabbitMQ struct {
	Conn *amqp.Connection

Ch *amqp.Channel
}


func NewRabbitMQ(url string) (*RabbitMQ, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()

		return nil, fmt.Errorf("failed to create RabbitMQ channel: %w", err)
	}

	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}, nil
}

func(r *RabbitMQ) DeclareQueue(name string ) error{
	_,err:= r.Ch.QueueDeclare(name , true,
		false,
		false,
		false,
		nil,)
		return err
}
func (r *RabbitMQ) Publish(
	queueName string,
	message []byte,
) error {

	return r.Ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	)
}