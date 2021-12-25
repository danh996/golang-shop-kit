package kafka_utils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	msg = &Message{
		Name: "Mr. Quang Teo",
		Age:  28,
	}
	topic = &Topic{
		Name: "simple_rule_test",
		ID:   uuid.New(),
	}
	Address            = "localhost:9092"
	PublishClientID    = "publisher_test"
	SubscriberClientID = "subscriber_test"
	Partition          = int32(0)
)

func Test_Kafka(t *testing.T) {
	p, err := NewPublisher(PublishClientID, Address)
	assert.NoError(t, err)

	err = p.Publish(topic, nil, msg)
	assert.NoError(t, err)

	c, err := NewSubscriber(topic, SubscriberClientID, Address, Partition)
	assert.NoError(t, err)

	fn := func(msg []byte) {
		var me *Message
		fmt.Println(string(msg))
		err := json.Unmarshal(msg, &me)
		// fmt.Printf("%+v", err)

		assert.NoError(t, err)
		assert.Equal(t, me, msg)
		fmt.Println(me)
	}
	c.Subscribe(fn)
}
