/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/9/6 6:38 PM
 * @desc: about the role of class.
 */

package kafkas

import "github.com/IBM/sarama"

type AdminClient struct {
	sarama.ClusterAdmin
}

type TopicDetail struct {
	sarama.TopicDetail
}

type MatchingAcl struct {
	sarama.MatchingAcl
}

type Acl struct {
	sarama.Acl
}

type Resource struct {
	sarama.Resource
}

type GroupDescription struct {
	sarama.GroupDescription
}

type LeaveGroupResponse struct {
	sarama.LeaveGroupResponse
}

type OffsetFetchResponse struct {
	sarama.OffsetFetchResponse
}

type PartitionReplicaReassignmentsStatus struct {
	sarama.PartitionReplicaReassignmentsStatus
}

func InitKafka(addrs []string, offsetOldest, isSync, randomPart bool, retryMax int) (*AdminClient, error) {
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
	clusterAdmin, err := sarama.NewClusterAdmin(addrs, config)
	return &AdminClient{clusterAdmin}, err
}

func InitKafkaPlain(addrs []string, username, password string, offsetOldest, isSync, randomPart bool, retryMax int) (*AdminClient, error) {
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
	clusterAdmin, err := sarama.NewClusterAdmin(addrs, config)
	return &AdminClient{clusterAdmin}, err
}

func InitKafkaScram(addrs []string, username, password string, offsetOldest, isSync, randomPart bool, retryMax int) (*AdminClient, error) {
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
	clusterAdmin, err := sarama.NewClusterAdmin(addrs, config)
	return &AdminClient{clusterAdmin}, err
}

func (c *AdminClient) TopicCreate(topic string, partitions int32, replicationFactor int16) error {
	topicConfig := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries:     map[string]*string{},
	}
	err := c.CreateTopic(topic, topicConfig, false)
	return err
}

func (c *AdminClient) TopicDelete(topic string) error {
	err := c.DeleteTopic(topic)
	return err
}

func (c *AdminClient) TopicList() (map[string]TopicDetail, error) {
	topicList, err := c.ListTopics()
	result := map[string]TopicDetail{}
	for k, v := range topicList {
		result[k] = TopicDetail{v}
	}
	return result, err
}

func (c *AdminClient) ACLCreate(resource Resource, acls []Acl) error {
	var aclList []*sarama.Acl
	for _, acl := range acls {
		aclList = append(aclList, &acl.Acl)
	}
	aclResource := &sarama.ResourceAcls{
		Resource: resource.Resource,
		Acls:     aclList,
	}
	err := c.CreateACLs([]*sarama.ResourceAcls{aclResource})
	return err
}

func (c *AdminClient) ACLDelete(acl string) ([]MatchingAcl, error) {
	aclFilter := sarama.AclFilter{
		ResourceName: &acl,
	}
	matchList, err := c.DeleteACL(aclFilter, false)
	var result []MatchingAcl
	for _, m := range matchList {
		result = append(result, MatchingAcl{m})
	}
	return result, err
}

func (c *AdminClient) ACLList() (map[string]TopicDetail, error) {
	topicList, err := c.ListTopics()
	result := map[string]TopicDetail{}
	for k, v := range topicList {
		result[k] = TopicDetail{v}
	}
	return result, err
}

func (c *AdminClient) PartitionCreate(topic string, count int32, assignment [][]int32) error {
	err := c.CreatePartitions(topic, count, assignment, false)
	return err
}

func (c *AdminClient) PartitionReassignAlter(topic string, assignment [][]int32) error {
	err := c.AlterPartitionReassignments(topic, assignment)
	return err
}

func (c *AdminClient) PartitionReassignList(topics string, partitions []int32) (map[string]map[int32]*PartitionReplicaReassignmentsStatus, error) {
	statusList, err := c.ListPartitionReassignments(topics, partitions)
	var topicStatus map[string]map[int32]*PartitionReplicaReassignmentsStatus
	for s, m := range statusList {
		status := map[int32]*PartitionReplicaReassignmentsStatus{}
		for k, v := range m {
			status[k] = &PartitionReplicaReassignmentsStatus{*v}
		}
		topicStatus[s] = status
	}
	return topicStatus, err
}

func (c *AdminClient) GroupDelete(group string) error {
	err := c.DeleteConsumerGroup(group)
	return err
}

func (c *AdminClient) GroupOffsetDelete(group string, topic string, partition int32) error {
	err := c.DeleteConsumerGroupOffset(group, topic, partition)
	return err
}

func (c *AdminClient) GroupList() (map[string]string, error) {
	result, err := c.ListConsumerGroups()
	return result, err
}

func (c *AdminClient) GroupOffsetsList(group string, topicPartitions map[string][]int32) (*OffsetFetchResponse, error) {
	response, err := c.ListConsumerGroupOffsets(group, topicPartitions)
	return &OffsetFetchResponse{*response}, err
}

func (c *AdminClient) GroupRemoveMember(groupId string, groupInstanceIds []string) (*LeaveGroupResponse, error) {
	response, err := c.RemoveMemberFromConsumerGroup(groupId, groupInstanceIds)
	return &LeaveGroupResponse{*response}, err
}

func (c *AdminClient) GroupDescribe(groups []string) ([]*GroupDescription, error) {
	groupList, err := c.DescribeConsumerGroups(groups)
	var result []*GroupDescription
	for _, g := range groupList {
		result = append(result, &GroupDescription{*g})
	}
	return result, err
}
