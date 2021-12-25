package kafka_utils

import (
	"fmt"

	"gitlab.com/canco1/canco-kit/logger"
)

type Subscriber interface {
	Subscribe(fn func(msg []byte))
}

type subscriber struct {
	Address   string
	Consumer  *Consumer
	Partition int32
	Topic     *Topic
	logger    logger.Logger
}

func (s *subscriber) Subscribe(fn func(msg []byte)) {
	for {
		select {
		case msg := <-s.Consumer.Consumer.Messages():
			fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			fn(msg.Value)
		case err := <-s.Consumer.Consumer.Errors():
			// s.logger.Error("consumer got error", zap.Error(err))
			fmt.Println(err)
		}
	}
}

func NewSubscriber(topic *Topic, clientID, address string, partition int32) (Subscriber, error) {

	consumer, err := NewConsumer(clientID, address, partition, topic)
	if err != nil {
		return nil, err
	}

	return &subscriber{
		Consumer:  consumer,
		Address:   address,
		Partition: partition,
		Topic:     topic,
	}, nil
}
