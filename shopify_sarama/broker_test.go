package shopify_sarama

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBroker(t *testing.T) {
	broker := sarama.NewBroker(brokersUrl[0])
	err := broker.Open(nil)
	require.NoError(t, err)
	defer func() {
		if err := broker.Close(); err != nil {
			t.Error(err)
		}
	}()

	request := sarama.MetadataRequest{Topics: topicList}
	response, err := broker.GetMetadata(&request)
	require.NoError(t, err)

	fmt.Printf("Custer Id(%d), Controller Id(%d)\n", response.ClusterID, response.ControllerID)
	fmt.Printf("There are %dbrokers active in the clusters.\n", len(response.Brokers))
	for _, broker := range response.Brokers {
		fmt.Printf("\t- Broker Id: %d\n", broker.ID())
	}
	fmt.Printf("There are %dtopics active in the clusters.\n", len(response.Topics))
	for _, topic := range response.Topics {
		fmt.Printf("\t- Topic(%s) has %dpartitions\n", topic.Name, len(topic.Partitions))
		for i, partition := range topic.Partitions {
			fmt.Printf("\t\t- Partition%d: id(%d), leader(%d), isr(%v), replicas(%v), offline replicas(%v)\n",
				i+1,
				partition.ID,
				partition.Leader,
				partition.Isr,
				partition.Replicas,
				partition.OfflineReplicas,
			)
		}
	}
}
