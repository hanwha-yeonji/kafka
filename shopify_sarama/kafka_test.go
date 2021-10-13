package sarama_test

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/require"
	"testing"
)

var msgMap = map[string][]string{
	"test": {"value1", "value2", "value3"},
}
var topicList = []string{
	"test",
}
var brokersUrl = []string{"localhost:19092", "localhost:29092", "localhost:39092"}

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
			require.True(t, offset == int64(i+1) && offset == msg.Offset, "the first msg should have been assigned offset ", i+1, ", but got ", msg.Offset)
			fmt.Printf("Message is stored in topic(%s) / partition(%d) / offset(%d)\n", topic, partition, offset)
		}
	}
}

func TestConsumer(t *testing.T) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	worker, err := sarama.NewConsumer(brokersUrl, config)
	require.NoError(t, err)

	cs, err := worker.ConsumePartition(topicList[0], 0, sarama.OffsetOldest)
	require.NoError(t, err)

	fmt.Println(cs)
}
