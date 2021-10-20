package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
)

type topicConsumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewTopicConsumer(brokersUrl []string, topic string) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return &topicConsumer{consumer: consumer, topic: topic}, nil
}

func (c *topicConsumer) Close() error {
	return c.consumer.Close()
}

func (c *topicConsumer) Consume(msgHandler MsgHandler, errHandler ErrHandler) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	msgCount := 0

	msgConsumer, consumerErr, err := c.consume()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgConsumer:
				msgCount++
				msgHandler(msg)
			case err := <-consumerErr:
				msgCount++
				errHandler(err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Printf("Topic(%s) Interrupt is detected\n", c.topic)
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Printf("Processed %dmessages\n", msgCount)

	return nil
}

func (c *topicConsumer) consume() (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError, error) {
	msgChan := make(chan *sarama.ConsumerMessage)
	errChan := make(chan *sarama.ConsumerError)

	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Start consuming topic: ", c.topic)
	for _, partition := range partitions {
		partitionConsumer, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetOldest)
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("Start consuming topic(%s): partition%d\n", c.topic, partition)
		go func() {
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					msgChan <- msg
				case err := <-partitionConsumer.Errors():
					errChan <- err
				}
			}
		}()
	}

	return msgChan, errChan, nil
}
