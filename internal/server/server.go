package server

import (
	"github.com/streadway/amqp"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
	"log"
	"context"
	"encoding/json"
)

type Server struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue

	audit *AuditServer
}

func NewServer(amqpURL, queueName string, auditService AuditService) (*Server, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"logs",  // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	auditServer := NewAuditServer(auditService)

	return &Server{
		conn:  conn,
		ch:    ch,
		queue: q,
		audit: auditServer,
	}, nil
}

func (s *Server) ConsumeMessages(ctx context.Context) {
	msgs, err := s.ch.Consume(
		s.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("failed to register a consumer: ", err)
	}

	go func() {
		for d := range msgs {
			var logMessage audit.LogMessage
			if err := json.Unmarshal(d.Body, &logMessage); err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}
			
			ctx := context.WithValue(context.Background(), "metadata", logMessage.Context)
			log.Printf("Received a message: %+v", logMessage.LogItem)
			err := s.audit.Insert(ctx, &logMessage.LogItem)
			if err != nil {
				log.Printf("failed to log message: %v", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

func (s *Server) Close() {
	s.ch.Close()
	s.conn.Close()
}
