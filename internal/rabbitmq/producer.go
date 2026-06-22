package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewProducer(url string) (*Producer, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Pastikan exchange dibuat
	err = ch.ExchangeDeclare(
		"fleet.events", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Producer{conn: conn, ch: ch}, nil
}

func (p *Producer) PublishGeofenceEvent(ctx context.Context, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(ctx,
		"fleet.events",   // exchange
		"geofence.entry", // routing key
		false,            // mandatory
		false,            // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

func (p *Producer) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
