/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/8/9 20:08
 * @desc: about the role of class.
 */

package rabbitmq

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Rabbitmq 初始化rabbitmq连接
type Rabbitmq struct {
	Conn      *Connection //连接
	Type      string      `json:"type"`       //消息类型
	Err       error       `json:"err"`        //错误信息
	Channel   *Channel    `json:"channel"`    //管道
	QueueName string      `json:"queue_name"` //队列名称
	Exchange  string      `json:"exchange"`   //交换机
	RouteKey  string      `json:"route_key"`  // 路由名称
	Key       string      `json:"key"`        //key Simple模式 几乎用不到
	MqUrl     string      `json:"mq_url"`     //连接信息
}

func New(username, password, ip string, port int64) (*Rabbitmq, error) {
	if len(username) < 1 {
		username = "guest"
	}
	if len(password) < 1 {
		password = "guest"
	}
	if len(ip) < 1 {
		ip = "localhost"
	}
	if port < 1 {
		port = 5672
	}
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, ip, port)
	conn, err := Dial(url)
	if err != nil {
		return nil, err
	}
	rabbitmq := &Rabbitmq{
		Conn: conn,
	}
	return rabbitmq, nil
}

func NewByCert(caFile, username, password, ip string, port int64) (*Rabbitmq, error) {

	// 加载CA证书 "ca.pem"
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cfg := &tls.Config{
		RootCAs: caCertPool,
	}

	if len(username) < 1 {
		username = "guest"
	}
	if len(password) < 1 {
		password = "guest"
	}
	if len(ip) < 1 {
		ip = "localhost"
	}
	if port < 1 {
		port = 5672
	}
	url := fmt.Sprintf("amqps://%s:%s@%s:%d/", username, password, ip, port)
	conn, err := DialTLS(url, cfg)
	if err != nil {
		return nil, err
	}
	rabbitmq := &Rabbitmq{
		Conn: conn,
	}
	return rabbitmq, nil
}

func (r *Rabbitmq) CreateQueue(name string) error {
	ch, err := r.Conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (r *Rabbitmq) PublishQueue(exchange, key string, body interface{}) error {
	bys, err := json.Marshal(body)
	if err != nil {
		return err
	}
	ch, err := r.Conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,
		Publishing{
			DeliveryMode: Persistent,
			ContentType:  "text/plain",
			Body:         bys,
		})
	return err
}

func (r *Rabbitmq) GetReadyCount(name string) (int, error) {
	count := 0
	ch, err := r.Conn.Channel()
	defer ch.Close()
	if err != nil {
		return count, err
	}
	state, err := ch.QueueInspect(name)
	if err != nil {
		return count, err
	}
	return state.Messages, nil
}

func (r *Rabbitmq) GetConsumCount(name string) (int, error) {
	count := 0
	ch, err := r.Conn.Channel()
	defer ch.Close()
	if err != nil {
		return count, err
	}
	state, err := ch.QueueInspect(name)
	if err != nil {
		return count, err
	}
	return state.Consumers, nil
}

func (r *Rabbitmq) ClearQueue(name string) error {
	ch, err := r.Conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}
	_, err = ch.QueuePurge(name, false)
	if err != nil {
		return err
	}
	return nil
}

// 创建RabbitMQ结构体实例
func NewByStruct(queuename, exchange, key, mqUrl string) (*Rabbitmq, error) {
	rabbitmq := &Rabbitmq{QueueName: queuename, Exchange: exchange, Key: key, MqUrl: mqUrl}
	var err error
	//创建rabbitmq连接
	rabbitmq.Conn, err = Dial(rabbitmq.MqUrl)
	if err != nil {
		return nil, err
	}

	rabbitmq.Channel, err = rabbitmq.Conn.Channel()
	return rabbitmq, nil
}

// 断开channel和connection
func (r *Rabbitmq) Destory() {
	r.Channel.Close()
	r.Conn.Close()
}

// 订阅模式生成
func (r *Rabbitmq) PublishPub(message interface{}) error {
	bys, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//尝试创建交换机，不存在创建
	err = r.Channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"fanout",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//2 发送消息
	err = r.Channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: bys,
		})
	return err
}

// 订阅模式消费端代码
func (r *Rabbitmq) RecieveSub(fn func(msg []byte)) error {
	//尝试创建交换机，不存在创建
	err := r.Channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"fanout",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//2试探性创建队列，创建队列
	q, err := r.Channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//绑定队列到exchange中
	err = r.Channel.QueueBind(
		q.Name,
		"", //在pub/sub模式下，这里的key要为空
		r.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//消费消息
	msgs, err := r.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fn(d.Body)
		}
	}()

	<-forever
	return nil
}

// 话题模式发送信息
func (r *Rabbitmq) PublishTopic(message interface{}) error {
	bys, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//尝试创建交换机，不存在创建
	err = r.Channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 话题模式
		"topic",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//2发送信息
	err = r.Channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: bys,
		})
	return err
}

// 话题模式接收信息
// 要注意key
// 其中* 用于匹配一个单词，#用于匹配多个单词（可以是零个）
// 匹配 表示匹配imooc.* 表示匹配imooc.hello,但是imooc.hello.one需要用imooc.#才能匹配到
func (r *Rabbitmq) RecieveTopic(fn func(msg []byte)) error {
	//尝试创建交换机，不存在创建
	err := r.Channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 话题模式
		"topic",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//2试探性创建队列，创建队列
	q, err := r.Channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//绑定队列到exchange中
	err = r.Channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//消费消息
	msgs, err := r.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fn(d.Body)
		}
	}()
	<-forever
	return nil
}

// 路由模式发送信息
func (r *Rabbitmq) PublishRouting(message interface{}) error {
	bys, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//尝试创建交换机，不存在创建
	err = r.Channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"direct",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//发送信息
	err = r.Channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: bys,
		})
	return err
}

// 路由模式接收信息
func (r *Rabbitmq) RecieveRouting(fn func(msg []byte)) error {
	//尝试创建交换机，不存在创建
	err := r.Channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"direct",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//2试探性创建队列，创建队列
	q, err := r.Channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//绑定队列到exchange中
	err = r.Channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//消费消息
	msgs, err := r.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fn(d.Body)
		}
	}()

	<-forever
	return nil
}

// 简单模式Step:2、简单模式下生产代码
func (r *Rabbitmq) PublishSimple(message interface{}) error {
	bys, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err = r.Channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,
		//是否具有排他性 true表示自己可见 其他用户不能访问
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		//额外数据
		nil,
	)
	if err != nil {
		return err
	}

	//2.发送消息到队列中
	err = r.Channel.Publish(
		//默认的Exchange交换机是default,类型是direct直接类型
		r.Exchange,
		//要赋值的队列名称
		r.QueueName,
		//如果为true，根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息还给发送者
		false,
		//消息
		Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: bys,
		})
	return err
}

func (r *Rabbitmq) ConsumeSimple(fn func(msg []byte)) error {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err := r.Channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外数据
		nil,
	)
	if err != nil {
		return err
	}
	//接收消息
	msgs, err := r.Channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true,表示不能同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		return err
	}
	forever := make(chan bool)

	//启用协程处理
	go func() {
		for d := range msgs {
			fn(d.Body)
		}
	}()

	<-forever
	return nil
}

func (r *Rabbitmq) ConsumeWorker(consumerName string, fn func(msg []byte)) error {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err := r.Channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外数据
		nil,
	)
	if err != nil {
		return err
	}
	//接收消息
	msgs, err := r.Channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		consumerName,
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true,表示不能同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		return err
	}
	forever := make(chan bool)

	//启用协程处理
	go func() {
		for d := range msgs {
			fn(d.Body)
		}
	}()

	<-forever
	return nil
}

// 获取到交换机
func (r *Rabbitmq) getExchange(exchange string) error {
	err := r.Channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable 持久化消息
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	return err
}

// 发送消息
func (r *Rabbitmq) SendMessageByKey(key string, msg interface{}) error {
	bys, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// 推送消息
	err = r.Channel.Publish(
		r.Exchange, // exchange
		key,        // routing key
		false,      // mandatory
		false,      // immediate
		Publishing{
			ContentType: "text/plain",
			Body:        bys,
		})
	return err
}

// 获取到消息
func (r *Rabbitmq) GetMessage(key string, fn func(msg []byte)) error {
	// 存储临时交换队列
	q, err := r.Channel.QueueDeclare(
		r.QueueName, // name
		true,        // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	err = r.Channel.QueueBind(
		q.Name,     // queue name
		key,        // routing key
		r.Exchange, // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	// 设置逐个消费消息
	err = r.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := r.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	ret := make(chan string)
	go func() {
		for {
			select {
			case d := <-msgs:
				fn(d.Body)
				d.Ack(false) // 标记消息被消费掉了
			}
		}
	}()

	<-ret
	return nil
}

// NewRabbitMQ 新建 rabbitmq 实例
func NewByExchange(username, password, ip string, port int64, vhost, exchange, route, queue string) (*Rabbitmq, error) {

	// 建立amqp链接
	conn, err := Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d%s",
		username,
		password,
		ip,
		port,
		"/"+strings.TrimPrefix(vhost, "/"),
	))
	if err != nil {
		return nil, err
	}
	r := &Rabbitmq{Conn: conn}

	// 建立channel通道
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	r.Channel = ch

	// 声明exchange交换器
	err = r.declareExchange(exchange, nil)
	return r, err
}

// SendMessage 发送普通消息
func (r *Rabbitmq) SendMessage(message interface{}) error {
	bys, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = r.Channel.Publish(
		r.Exchange, // exchange
		r.RouteKey, // route key
		false,
		false,
		Publishing{
			ContentType: "text/plain",
			Body:        bys,
		},
	)
	return err
}

// SendDelayMessage 发送延迟消息
func (r *Rabbitmq) SendDelayMessage(message interface{}, delayTime int) error {
	delayQueueName := r.QueueName + "_delay:" + strconv.Itoa(delayTime)
	delayRouteKey := r.RouteKey + "_delay:" + strconv.Itoa(delayTime)

	bys, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 定义延迟队列(死信队列)
	dq, err := r.declareQueue(
		delayQueueName,
		Table{
			"x-dead-letter-exchange":    r.Exchange, // 指定死信交换机
			"x-dead-letter-routing-key": r.RouteKey, // 指定死信routing-key
		},
	)
	if err != nil {
		return err
	}

	// 延迟队列绑定到exchange
	err = r.bindQueue(dq.Name, delayRouteKey, r.Exchange)
	if err != nil {
		return err
	}

	// 发送消息，将消息发送到延迟队列，到期后自动路由到正常队列中
	err = r.Channel.Publish(
		r.Exchange,
		delayRouteKey,
		false,
		false,
		Publishing{
			ContentType: "text/plain",
			Body:        bys,
			Expiration:  strconv.Itoa(delayTime * 1000),
		},
	)
	return err
}

// Consume 获取消费消息
func (r *Rabbitmq) Consume(fn func(msg []byte)) error {
	// 声明队列
	q, err := r.declareQueue(r.QueueName, nil)
	if err != nil {
		return err
	}
	// 队列绑定到exchange
	err = r.bindQueue(q.Name, r.RouteKey, r.Exchange)
	if err != nil {
		return err
	}
	// 设置Qos
	err = r.Channel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	// 监听消息
	msgs, err := r.Channel.Consume(
		q.Name, // queue name,
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	forever := make(chan bool) //注册在主进程，不需要阻塞

	go func() {
		for d := range msgs {
			fn(d.Body)
			d.Ack(false)
		}
	}()

	<-forever
	return nil
}

// Close 关闭链接
func (r *Rabbitmq) Close() {
	r.Channel.Close()
	r.Conn.Close()
}

// declareQueue 定义队列
func (r *Rabbitmq) declareQueue(name string, args Table) (Queue, error) {
	q, err := r.Channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)
	return q, err
}

// declareQueue 定义交换器
func (r *Rabbitmq) declareExchange(exchange string, args Table) error {
	err := r.Channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		args,
	)
	return err
}

// bindQueue 绑定队列
func (r *Rabbitmq) bindQueue(queue, routekey, exchange string) error {
	err := r.Channel.QueueBind(
		queue,
		routekey,
		exchange,
		false,
		nil,
	)
	return err
}
