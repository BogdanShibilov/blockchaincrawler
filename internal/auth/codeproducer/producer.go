package codeproducer

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/config"
)

type CodeProducer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

func NewCodeProducer(cfg config.CodeProducer) (*CodeProducer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create new async producer: %w", err)
	}

	return &CodeProducer{
		asyncProducer: asyncProducer,
		topic:         cfg.Topic,
	}, nil
}

func (p *CodeProducer) ProduceCode(code *ConfirmationCode) error {
	codeJSON, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("failed to marshal confirmation code: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(codeJSON),
	}

	p.asyncProducer.Input() <- msg
	return nil
}

func (p *CodeProducer) Errors() <-chan *sarama.ProducerError {
	return p.asyncProducer.Errors()
}
