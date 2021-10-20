package consumer

import "github.com/Shopify/sarama"

type MsgHandler func(*sarama.ConsumerMessage)
type ErrHandler func(*sarama.ConsumerError)

type Consumer interface {
	Consume(msgHandler MsgHandler, errHandler ErrHandler) error
	Close() error
}
