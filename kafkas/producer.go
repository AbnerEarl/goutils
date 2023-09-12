/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/9/6 6:38 PM
 * @desc: about the role of class.
 */

package kafkas

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type ProducerClient struct {
	sarama.SyncProducer
}

type Message struct {
	sarama.ProducerMessage
}

type StringEncoder sarama.StringEncoder

func (s StringEncoder) Encode() ([]byte, error) {
	return []byte(s), nil
}

func (s StringEncoder) Length() int {
	return len(s)
}

type ByteEncoder sarama.ByteEncoder

func (b ByteEncoder) Encode() ([]byte, error) {
	return b, nil
}

func (b ByteEncoder) Length() int {
	return len(b)
}

func InitProducer(addrs []string, isSync, randomPart bool, retryMax int) (*ProducerClient, error) {
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
	syncProducer, err := sarama.NewSyncProducer(addrs, config)
	return &ProducerClient{syncProducer}, err
}

func InitProducerPlain(addrs []string, username, password string, isSync, randomPart bool, retryMax int) (*ProducerClient, error) {
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
	syncProducer, err := sarama.NewSyncProducer(addrs, config)
	return &ProducerClient{syncProducer}, err
}

func InitProducerScram(addrs []string, username, password string, isSync, randomPart bool, retryMax int) (*ProducerClient, error) {
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
	syncProducer, err := sarama.NewSyncProducer(addrs, config)
	return &ProducerClient{syncProducer}, err
}

func (c *ProducerClient) ProducerMessage(topic string, msg interface{}) (partition int32, offset int64, err error) {
	bys, _ := json.Marshal(msg)
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(bys),
	}
	return c.SendMessage(message)
}

func (c *ProducerClient) ProducerMessages(topic string, msgs []interface{}) (err error) {
	var messages []*sarama.ProducerMessage
	for msg := range msgs {
		bys, _ := json.Marshal(msg)
		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(bys),
		}
		messages = append(messages, message)
	}
	return c.SendMessages(messages)
}

func (c *ProducerClient) ProducerCustom(msg Message) (partition int32, offset int64, err error) {
	return c.SendMessage(&msg.ProducerMessage)
}

func (c *ProducerClient) ProducerCustoms(msgs []Message) (err error) {
	var messages []*sarama.ProducerMessage
	for _, msg := range msgs {
		messages = append(messages, &msg.ProducerMessage)
	}
	return c.SendMessages(messages)
}
