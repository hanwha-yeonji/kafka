package shopify_sarama

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hanwha-yeonji/kafka/shopify_sarama/producer"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

const sendMsgTimes = 5

var producerMsgHandler = func(msg *sarama.ProducerMessage) {
	fmt.Printf("Topic(%s) Successfully send a messages. %dMessage(key/value): (%s/%s)\n", msg.Topic, msg.Offset, msg.Key, msg.Value)
}
var producerErrHandler = func(err *sarama.ProducerError) {
	if err == nil {
		fmt.Printf("Failed send a message. Err message: %v\n", err)
	} else {
		fmt.Printf("Failed send a message. Send message %v, Err message: %v\n", err.Msg, err.Err)
	}
}

func TestGeneralSyncProducerSendMessageNoKey(t *testing.T) {
	producer, err := producer.NewGeneralSyncProducer(brokersUrl)
	require.NoError(t, err)

	defer func() {
		if err := producer.Close(); err != nil {
			t.Error(err)
		}
	}()

	for _, topic := range topicList {
		for i := 0; i < sendMsgTimes; i++ {
			msgTxt := generateRandomString(10)
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(msgTxt),
			}
			fmt.Println(msg)
			err := producer.SendMessage(msg)
			require.NoError(t, err)
		}
	}
}

func TestGeneralSyncProducerSendMessageWithKey(t *testing.T) {
	producer, err := producer.NewGeneralSyncProducer(brokersUrl)
	require.NoError(t, err)

	defer func() {
		if err := producer.Close(); err != nil {
			t.Error(err)
		}
	}()

	keys := []string{"key1", "key2"}

	for i := 0; i < sendMsgTimes; i++ {
		msgTxt := generateRandomString(10)
		msg := &sarama.ProducerMessage{
			Topic: topicList[1],
			Key:   sarama.StringEncoder(keys[i%2]),
			Value: sarama.StringEncoder(msgTxt),
		}
		err := producer.SendMessage(msg)
		require.NoError(t, err)
	}
}

func TestTopicSyncProducerSendMessageWithKey(t *testing.T) {
	producer, err := producer.NewTopicSyncProducer(brokersUrl, topicList[0])
	require.NoError(t, err)

	defer func() {
		if err := producer.Close(); err != nil {
			t.Error(err)
		}
	}()

	keys := []string{"key1", "key2"}

	for i := 0; i < sendMsgTimes; i++ {
		msgTxt := generateRandomString(10)
		msg := &sarama.ProducerMessage{
			Key:   sarama.StringEncoder(keys[i%2]),
			Value: sarama.StringEncoder(msgTxt),
		}
		err := producer.SendMessage(msg)
		require.NoError(t, err)
	}
}

func TestTopicAsyncProducerSendMessageWithKey(t *testing.T) {
	producer, err := producer.NewTopicAsyncProducer(brokersUrl, topicList[0])
	require.NoError(t, err)

	defer func() {
		if err := producer.Close(); err != nil {
			t.Error(err)
		}
	}()

	go func() {
		producer.HandleMessage(producerMsgHandler, producerErrHandler)
	}()

	keys := []string{"key1", "key2"}

	for i := 0; i < sendMsgTimes; i++ {
		msgTxt := generateRandomString(10)
		msg := &sarama.ProducerMessage{
			Key:   sarama.StringEncoder(keys[i%2]),
			Value: sarama.StringEncoder(msgTxt),
		}
		err := producer.SendMessage(msg)
		require.NoError(t, err)
	}

	time.Sleep(time.Second * 5)
}

func randInt(max int) int {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return seed.Intn(max)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randInt(len(charset))]
	}

	return string(b)
}
