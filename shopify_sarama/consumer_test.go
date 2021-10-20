package shopify_sarama

import (
	"github.com/hanwha-yeonji/kafka/shopify_sarama/consumer"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGeneralConsumerConsume(t *testing.T) {
	consumer, err := consumer.NewGeneralConsumer(brokersUrl)
	require.NoError(t, err)

	defer func() {
		if err := consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	err = consumer.Consume(msgHandler, errHandler)
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

	err = consumer.Consume(msgHandler, errHandler)
	require.NoError(t, err)
}
