package rabbitmq

import (
	"Iris/datamodels"
	"Iris/services"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

// MQURL 连接信息 amqp://账号:密码@ip:host/vhost
const MQURL = "amqp://guest:guest@127.0.0.1:5672/"

// RabbitMQ rabbitMQ结构体
type RabbitMQ struct {
	conn      *amqp.Connection // 链接
	channel   *amqp.Channel    // 通道
	QueueName string           //队列名称
	Exchange  string           //交换机名称
	Key       string           //bind Key 名称
	Mqurl     string           //连接信息
	sync.Mutex
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// Destroy 断开 channel 和 connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close() // 断开 channel
	r.conn.Close()    // 断开 conn
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Printf("%s:%s", message, err)         // 打印错误
		panic(fmt.Sprintf("%s:%s", message, err)) // 抛出错误
	}
}

// NewRabbitMQSimple 创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	// todo 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	// todo 补上conn与channel
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishSimple 简单模式下队列生产
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	// todo 申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 首先放入名称
		false,       //是否持久化
		false,       //是否自动删除
		false,       //是否具有排他性
		false,       //是否阻塞处理
		nil,         //额外的属性
	)
	if err != nil {
		return err
	}
	//todo 调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange, // 此处为空
		r.QueueName,
		false, //如果为true，根据自身exchange类型和routeKey规则；无法找到符合条件的队列会把消息返还给发送者
		false, //如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return nil
}

// ConsumeSimple 简单模式下消费者
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	//todo 申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //是否自动删除
		false, //是否具有排他性
		false, //是否阻塞处理
		nil,   //额外的属性
	)
	if err != nil {
		fmt.Println(err)
	}

	//todo 消费者流控
	r.channel.Qos(
		1,     //当前消费者一次能接受的最大消息数量
		0,     //服务器传递的最大容量（以八位字节为单位）
		false, //如果设置为true 对channel可用
	)

	//todo 接收消息
	msg, err := r.channel.Consume(
		q.Name, // queue
		"",     //用来区分多个消费者 此处不区分
		false,  //是否自动应答
		false,  //是否独有
		false,  //设置为true，表示不能将同一个Connection中生产者发送的消息传递给这个Connection中的消费者
		false,  // 是否阻塞处理
		nil,    // 额外的属性
	)
	if err != nil {
		fmt.Println(err)
	}

	//todo 启用协程处理消息

	// 此处使用forever的意思为因为协程会始终监听消息(除非手动结束)
	// 手动结束才会进行 <-forever 有协程且一直尝试读取数据
	forever := make(chan bool)
	go func() {
		for d := range msg {
			// todo 消息逻辑处理
			message := &datamodels.Message{}
			err := json.Unmarshal(d.Body, message)
			if err != nil {
				fmt.Println(err)
			}
			// todo 插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}
			//todo 扣除商品数量
			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}
			d.Ack(false) //如果为true表示确认所有未确认的消息，为false表示确认当前消息
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
