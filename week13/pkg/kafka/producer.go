package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

var Producer *SyncProducer

type SyncProducer struct {
	syncProducer sarama.SyncProducer
}

func InitProducer(hosts []string) error {

	mqConfig := sarama.NewConfig()
	mqConfig.Producer.Return.Successes = true

	var err error

	if err = mqConfig.Validate(); err != nil {
		return fmt.Errorf("Kafka producer config invalidate. config: %v. err: %v", *mqConfig, err)
	}

	p, err := sarama.NewSyncProducer(hosts, mqConfig)
	if err != nil {
		return err
	}

	Producer = &SyncProducer{syncProducer: p}
	return nil
}

func (p *SyncProducer) Produce(topic string, key string, content string) error {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(content),
		Timestamp: time.Now(),
	}

	_, _, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Send msg failed: topic: %v, key: %v, content: %v, err: %v", topic, key, content, err)
	}

	return nil
}