package producer

import "github.com/Shopify/sarama"

type MsgHandler func(msgHandler *sarama.ProducerMessage)
type ErrHandler func(errHandler *sarama.ProducerError)

type SyncProducer interface {
	SendMessage(msg *sarama.ProducerMessage) error
	Close() error
}

type AsyncProducer interface {
	SendMessage(msg *sarama.ProducerMessage) error
	HandleMessage(msgHandler MsgHandler, errHandler ErrHandler)
	Close() error
}
