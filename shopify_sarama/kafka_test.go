package shopify_sarama

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var msgMap = map[string][]string{
	"test2": {"value1", "value2", "value3"},
}
var topicList = []string{
	"test", "test1", "test2",
}
var brokersUrl = []string{"localhost:19092", "localhost:29092", "localhost:39092"}

var msgHandler = func(msg *sarama.ConsumerMessage) {
	fmt.Printf("Topic(%s) Received messages. %dMessage(key/value): (%s/%s)\n", msg.Topic, msg.Offset, msg.Key, msg.Value)
}
var errHandler = func(err *sarama.ConsumerError) {
	fmt.Printf("Topic(%s) Received error. Partition %d, Err message: %v\n", err.Topic, err.Partition, err.Err)
}
