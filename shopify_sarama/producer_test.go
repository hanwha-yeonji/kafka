package shopify_sarama

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	//config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	sp, err := sarama.NewSyncProducer(brokersUrl, config)
	require.NoError(t, err)
	defer func() {
		if err := sp.Close(); err != nil {
			t.Error(err)
		}
	}()

	for _, topic := range topicList {
		for i, msgTxt := range msgMap[topic] {
			msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(msgTxt)}
			fmt.Println(msg)
			// TODO: SendMessage timeout
			partition, offset, err := sp.SendMessage(msg)
			require.NoError(t, err)
			require.True(t, offset == int64(i) && offset == msg.Offset, fmt.Sprintf("the first msg should have been assigned offset %d, but got %d", i, msg.Offset))
			fmt.Printf("Message is stored in topic(%s) / partition(%d) / offset(%d)\n", topic, partition, offset)
		}
	}
}
