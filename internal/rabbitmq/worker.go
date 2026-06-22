package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewWorker(url string) (*Worker, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare("fleet.events", "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"geofence_alerts", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,           // queue name
		"geofence.entry", // routing key
		"fleet.events",   // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Worker{conn: conn, ch: ch}, nil
}

func (w *Worker) Start() {
	msgs, err := w.ch.Consume(
		"geofence_alerts", // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("RabbitMQ Worker started. Waiting for geofence alerts...")

	for d := range msgs {
		log.Printf("[RABBITMQ WORKER] Alert Diterima! Pesan: %s", string(d.Body))
	}
}

func (w *Worker) Close() {
	if w.ch != nil {
		w.ch.Close()
	}
	if w.conn != nil {
		w.conn.Close()
	}
}
