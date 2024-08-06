package utils

import (
	"fmt"

	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/svc"
	amqp "github.com/rabbitmq/amqp091-go"
)

const MainExchangeName = "main_exchange"

// RabbitMQTools provides methods for interacting with RabbitMQ, including sending messages.
type RabbitMQTools struct {
	conn         *amqp.Connection // RabbitMQ connection
	ch           *amqp.Channel    // RabbitMQ channel
	exchangeName string           // Name of the RabbitMQ exchange
}

// NewRabbitMQTools creates a new RabbitMQTools instance by establishing a channel and declaring an exchange.
func NewRabbitMQTools(svcCtx *svc.ServiceContext) (*RabbitMQTools, error) {
	// Create a new channel from the RabbitMQ connection
	ch, err := svcCtx.MQConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ channel: %w", err)
	}

	// Declare an exchange with the given name and type
	err = ch.ExchangeDeclare(
		MainExchangeName,
		amqp.ExchangeDirect, // Exchange type
		true,                // Durable: the exchange will survive server restarts
		false,               // Passive: if true, the exchange must already exist
		false,               // Internal: if true, only brokers can send to this exchange
		false,               // NoWait: if true, don't wait for server confirmation
		nil,                 // Arguments
	)
	if err != nil {
		ch.Close() // Close channel if exchange declaration fails
		return nil, fmt.Errorf("failed to declare RabbitMQ exchange: %w", err)
	}

	return &RabbitMQTools{
		conn:         svcCtx.MQConn,
		ch:           ch,
		exchangeName: MainExchangeName,
	}, nil
}

// SendMessage sends a message to a specific user and file queue.
// It declares a queue, binds it to the exchange with a routing key, and publishes the message.
func (rb *RabbitMQTools) SendMessage(userId int64, fileId string) error {
	routingKey := fmt.Sprintf("%d.%s", userId, fileId) // Create a routing key
	body := fmt.Sprintf("%d.%s", userId, fileId)       // Message body
	queueName := fmt.Sprintf("user_queue_%d", userId)  // Queue name for the user

	// Declare a queue with the given name
	_, err := rb.ch.QueueDeclare(
		queueName,
		true,  // Durable: the queue will survive server restarts
		false, // AutoDelete: if true, the queue is deleted when no longer in use
		false, // Exclusive: if true, the queue is only used by the current connection
		false, // NoWait: if true, don't wait for server confirmation
		nil,   // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare RabbitMQ queue: %w", err)
	}

	// Bind the queue to the exchange with the specified routing key
	err = rb.ch.QueueBind(queueName,
		routingKey,
		rb.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	// Publish the message to the exchange with the routing key
	err = rb.ch.Publish(rb.exchangeName,
		routingKey,
		false, // Mandatory: if true, the server will return an error if the message cannot be routed
		false, // Immediate: if true, the server will return an error if the message cannot be delivered immediately
		amqp.Publishing{
			ContentType: "text/plain", // Type of content being sent
			Body:        []byte(body), // Message body
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

// Close closes the RabbitMQ channel, releasing any resources associated with it.
func (rb *RabbitMQTools) Close() error {
	if err := rb.ch.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %w", err)
	}
	return nil
}
