package producer

import (
	"fmt"
	"github.com/Shopify/sarama"
)

type generalSyncProducer struct {
	producer sarama.SyncProducer
}

func NewGeneralSyncProducer(brokersUrl []string) (SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	//config.SyncProducer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	sp, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return &generalSyncProducer{producer: sp}, nil
}

func (p *generalSyncProducer) Close() error {
	return p.producer.Close()
}

func (p *generalSyncProducer) SendMessage(msg *sarama.ProducerMessage) error {
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s) / partition(%d) / offset(%d)\n", msg.Topic, partition, offset)
	return nil
}
