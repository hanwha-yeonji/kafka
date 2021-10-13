package segmentio_kafka_go

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var msgMap = map[string][]string{
	"test": {"value1", "value2", "value3"},
}
var topicList = []string{
	"test",
}
var brokersUrl = "127.0.0.1:9092"
var partition = 0

func TestWriteMessage(t *testing.T) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokersUrl, topicList[0], partition)
	require.NoError(t, err)

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	var msgs []kafka.Message
	for _, topic := range topicList {
		for _, msgTxt := range msgMap[topic] {
			msgs = append(msgs, kafka.Message{Value: []byte(msgTxt)})
		}
	}
	_, err = conn.WriteMessages(msgs...)
	require.NoError(t, err)

	defer func() {
		err := conn.Close()
		require.NoError(t, err)
	}()
}

func TestReadMessage(t *testing.T) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokersUrl, topicList[0], partition)
	require.NoError(t, err)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		require.NoError(t, err)

		fmt.Println(string(b))
	}

	defer func() {
		err := batch.Close()
		require.NoError(t, err)
	}()
}
