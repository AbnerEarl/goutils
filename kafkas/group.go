/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/9/6 6:38 PM
 * @desc: about the role of class.
 */

package kafkas

import (
	"context"
	"github.com/IBM/sarama"
)

type ConsumerGroupClient struct {
	sarama.ConsumerGroup
}

type ConsumerGroup struct {
	Size     int
	Messages [][]byte
}

func (c *ConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	tag := 0
	for msg := range claim.Messages() {
		c.Messages = append(c.Messages, msg.Value)
		tag++
		if c.Size > 0 && tag == c.Size {
			return nil
		}
	}
	return nil
}

func InitConsumerGroup(addrs []string, groupID string, isSync, offsetOldest, randomPart bool, retryMax int) (*ConsumerGroupClient, error) {
	/**
	 * @author: yangchangjia
	 * @email 1320259466@qq.com
	 * @date: 2023/9/7 10:08 AM
	 * @desc: about the role of function.
	 * @param addrs, the kafka cluster address, such as: []string{"localhost:9192","localhost:9292","localhost:9392"}
	 * @param username, the kafka username
	 * @param password, the kafka password
	 * @return null
	 */
	config := sarama.NewConfig()

	if isSync {
		config.Producer.RequiredAcks = sarama.WaitForAll
	}
	if randomPart {
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	} else {
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}
	if retryMax > 0 {
		config.Producer.Retry.Max = retryMax
	}
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	if offsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	consumer, err := sarama.NewConsumerGroup(addrs, groupID, config)
	return &ConsumerGroupClient{consumer}, err
}

func InitConsumerGroupScram(addrs []string, groupID, username, password string, isSync, offsetOldest, randomPart bool, retryMax int) (*ConsumerGroupClient, error) {
	/**
	 * @author: yangchangjia
	 * @email 1320259466@qq.com
	 * @date: 2023/9/7 10:08 AM
	 * @desc: about the role of function.
	 * @param addrs, the kafka cluster address, such as: []string{"localhost:9192","localhost:9292","localhost:9392"}
	 * @param username, the kafka username
	 * @param password, the kafka password
	 * @return null
	 */
	config := sarama.NewConfig()
	if len(username) > 0 {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: SHA256}
		}
	}

	if isSync {
		config.Producer.RequiredAcks = sarama.WaitForAll
	}
	if randomPart {
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	} else {
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}
	if retryMax > 0 {
		config.Producer.Retry.Max = retryMax
	}
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	if offsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	consumer, err := sarama.NewConsumerGroup(addrs, groupID, config)
	return &ConsumerGroupClient{consumer}, err
}
func InitConsumerGroupPlain(addrs []string, groupID, username, password string, isSync, offsetOldest, randomPart bool, retryMax int) (*ConsumerGroupClient, error) {
	/**
	 * @author: yangchangjia
	 * @email 1320259466@qq.com
	 * @date: 2023/9/7 10:08 AM
	 * @desc: about the role of function.
	 * @param addrs, the kafka cluster address, such as: []string{"localhost:9192","localhost:9292","localhost:9392"}
	 * @param username, the kafka username
	 * @param password, the kafka password
	 * @return null
	 */
	config := sarama.NewConfig()
	if len(username) > 0 {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}

	if isSync {
		config.Producer.RequiredAcks = sarama.WaitForAll
	}
	if randomPart {
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	} else {
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}
	if retryMax > 0 {
		config.Producer.Retry.Max = retryMax
	}
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	if offsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	consumer, err := sarama.NewConsumerGroup(addrs, groupID, config)
	return &ConsumerGroupClient{consumer}, err
}

func (c *ConsumerGroupClient) ConsumerMessage(topics []string) (value []byte, err error) {
	ctx := context.Background()
	handler := &ConsumerGroup{Size: 1}
	err = c.Consume(ctx, topics, handler)
	if err != nil {
		return nil, err
	}
	return handler.Messages[0], nil
}

func (c *ConsumerGroupClient) ConsumerMessages(topics []string, msgSize int) (value [][]byte, err error) {
	ctx := context.Background()
	handler := &ConsumerGroup{Size: msgSize}
	err = c.Consume(ctx, topics, handler)
	if err != nil {
		return nil, err
	}
	return handler.Messages, nil
}

func (c *ConsumerGroupClient) ConsumerCustom(topics []string, handler ConsumerGroup) (err error) {
	ctx := context.Background()
	err = c.Consume(ctx, topics, &handler)
	return err
}
