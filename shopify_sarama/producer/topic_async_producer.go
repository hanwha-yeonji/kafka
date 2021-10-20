package producer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
)

type topicAsyncProducer struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewTopicAsyncProducer(brokersUrl []string, topic string) (AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	//config.SyncProducer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	sp, err := sarama.NewAsyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return &topicAsyncProducer{producer: sp, topic: topic}, nil
}

func (p *topicAsyncProducer) Close() error {
	return p.producer.Close()
}

func (p *topicAsyncProducer) SendMessage(msg *sarama.ProducerMessage) error {
	msg.Topic = p.topic

	p.producer.Input() <- msg

	return nil
}

func (p *topicAsyncProducer) HandleMessage(msgHandler MsgHandler, errHandler ErrHandler) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case msg := <-p.producer.Successes():
				msgHandler(msg)
			case err := <-p.producer.Errors():
				errHandler(err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
}
