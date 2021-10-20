package producer

import (
	"fmt"
	"github.com/Shopify/sarama"
)

type topicSyncProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewTopicSyncProducer(brokersUrl []string, topic string) (SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	//config.SyncProducer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	sp, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return &topicSyncProducer{producer: sp, topic: topic}, nil
}

func (p *topicSyncProducer) Close() error {
	return p.producer.Close()
}

func (p *topicSyncProducer) SendMessage(msg *sarama.ProducerMessage) error {
	msg.Topic = p.topic
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s) / partition(%d) / offset(%d)\n", msg.Topic, partition, offset)
	return nil
}
