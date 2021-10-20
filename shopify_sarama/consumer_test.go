package shopify_sarama

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hanwha-yeonji/kafka/shopify_sarama/consumer"
	"github.com/stretchr/testify/require"
	"testing"
)

var consumerMsgHandler = func(msg *sarama.ConsumerMessage) {
	fmt.Printf("Topic(%s) Received messages. %dMessage(key/value): (%s/%s)\n", msg.Topic, msg.Offset, msg.Key, msg.Value)
}
var consumerErrHandler = func(err *sarama.ConsumerError) {
	fmt.Printf("Topic(%s) Received error. Partition %d, Err message: %v\n", err.Topic, err.Partition, err.Err)
}

func TestGeneralConsumerConsume(t *testing.T) {
	consumer, err := consumer.NewGeneralConsumer(brokersUrl)
	require.NoError(t, err)

	defer func() {
		if err := consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	err = consumer.Consume(consumerMsgHandler, consumerErrHandler)
	require.NoError(t, err)
}

func TestTopicConsumerConsume(t *testing.T) {
	topic := topicList[0]
	consumer, err := consumer.NewTopicConsumer(brokersUrl, topic)
	require.NoError(t, err)

	defer func() {
		if err := consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	err = consumer.Consume(consumerMsgHandler, consumerErrHandler)
	require.NoError(t, err)
}
