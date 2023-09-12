/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/9/6 6:38 PM
 * @desc: about the role of class.
 */

package kafkas

import (
	"github.com/IBM/sarama"
	"sync"
)

type ConsumerClient struct {
	sarama.Consumer
	Offset int64
}

func InitConsumer(addrs []string, offsetOldest, isSync, randomPart bool, retryMax int) (*ConsumerClient, error) {
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
	consumer, err := sarama.NewConsumer(addrs, config)
	return &ConsumerClient{consumer, config.Consumer.Offsets.Initial}, err
}

func InitConsumerPlain(addrs []string, username, password string, offsetOldest, isSync, randomPart bool, retryMax int) (*ConsumerClient, error) {
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
	consumer, err := sarama.NewConsumer(addrs, config)
	return &ConsumerClient{consumer, config.Consumer.Offsets.Initial}, err
}

func InitConsumerScram(addrs []string, username, password string, offsetOldest, isSync, randomPart bool, retryMax int) (*ConsumerClient, error) {
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
	consumer, err := sarama.NewConsumer(addrs, config)
	return &ConsumerClient{consumer, config.Consumer.Offsets.Initial}, err
}

func (c *ConsumerClient) ConsumerMessage(topic string, partition int32) (value []byte, err error) {

	consumePartition, err := c.ConsumePartition(topic, partition, c.Offset)
	if err != nil {
		return nil, err
	}
	defer consumePartition.Close()
	msg := <-consumePartition.Messages()
	return msg.Value, nil
}

func (c *ConsumerClient) ConsumerMessages(topic string, fn func(message []byte) error) (err error) {
	partitionList, err := c.Partitions(topic)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(partitionList))
	for _, partition := range partitionList {
		consumePartition, err := c.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			return err
		}
		defer consumePartition.Close()
		go func(consumePartition sarama.PartitionConsumer) {
			for item := range consumePartition.Messages() {
				fn(item.Value)
			}
			wg.Done()
		}(consumePartition)
	}
	wg.Wait()
	return nil
}
