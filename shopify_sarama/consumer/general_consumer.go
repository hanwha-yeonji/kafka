package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
)

type generalConsumer struct {
	consumer sarama.Consumer
}

func NewGeneralConsumer(brokersUrl []string) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return &generalConsumer{consumer: consumer}, nil
}

func (c *generalConsumer) Close() error {
	return c.consumer.Close()
}

func (c *generalConsumer) Consume(msgHandler MsgHandler, errHandler ErrHandler) error {
	topics, _ := c.consumer.Topics()
	fmt.Println(len(topics), "topics")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	msgCount := 0

	for _, topic := range topics {
		msgConsumer, consumerErr, err := c.consume(topic)
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
					fmt.Printf("Topic(%s) Interrupt is detected\n", topic)
					doneCh <- struct{}{}
				}
			}
		}()

	}

	<-doneCh
	fmt.Printf("Processed %dmessages\n", msgCount)

	return nil
}

func (c *generalConsumer) consume(topic string) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError, error) {
	msgChan := make(chan *sarama.ConsumerMessage)
	errChan := make(chan *sarama.ConsumerError)

	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Start consuming topic: ", topic)
	for _, partition := range partitions {
		partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("Start consuming topic(%s): partition%d\n", topic, partition)
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
