package utils

import (
	"fmt"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/config"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/svc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/zrpc"
	"strconv"
	"strings"
	"testing"
	"time"
)

func NewRabbitMQConn(c config.Config) *amqp.Connection {
	var builder strings.Builder
	builder.WriteString(c.RabbitMQ.Protocol)
	builder.WriteString("://")
	builder.WriteString(c.RabbitMQ.Username)
	builder.WriteString(":")
	builder.WriteString(c.RabbitMQ.Password)
	builder.WriteString("@")
	builder.WriteString(c.RabbitMQ.Host)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(c.RabbitMQ.Port))
	dns := builder.String()

	conn, err := amqp.Dial(dns)
	if err != nil {
		panic("err dial to rabbitmq")
	}

	return conn
}

var rabbitmqConfig = config.Config{
	RpcServerConf: zrpc.RpcServerConf{},
	Datasource:    "",
	RabbitMQ: config.RabbitMQConf{
		Protocol: "amqp",
		Username: "guest",
		Password: "guest",
		Host:     "127.0.0.1",
		Port:     5672,
	},
}

var rabbitmqCtx = &svc.ServiceContext{
	Config: rabbitmqConfig,
	MQConn: NewRabbitMQConn(rabbitmqConfig),
}

func TestSendMessage(t *testing.T) {
	tool, err := NewRabbitMQTools(rabbitmqCtx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = tool.SendMessage(3746286, "48636524")
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("Send a message success")
}

func TestConsumeMessage(t *testing.T) {
	// TODO:Need be change by testing
	queue_name := fmt.Sprintf("user_queue_%d", 873526635)

	// 声明一个channel
	ch, err := rabbitmqCtx.MQConn.Channel()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		queue_name, // 队列名
		true,       // 是否持久化
		false,      // 是否自动删除
		false,      // 是否排他
		false,      // 其他参数
		nil,        // 额外参数
	)
	if err != nil {
		t.Fatalf("Failed to declare a queue: %v", err)
	}

	// 消费消息
	msgs, err := ch.Consume(
		q.Name, // 队列名
		"",     // 消费者标签
		true,   // 自动确认
		false,  // 是否排他
		false,  // 是否不再本地队列
		false,  // 是否无等待
		nil,    // 额外参数
	)
	if err != nil {
		t.Fatalf("Failed to register a consumer: %v", err)
	}

	// 创建消息通道来接收消息
	go func() {
		for msg := range msgs {
			t.Logf("Received a message: %s", msg.Body)
		}
	}()

	// 等待消息接收完成（这里可以添加超时逻辑以确保测试在合理时间内完成）
	// 这里简单的 sleep 用于模拟等待，你可以根据需要调整
	select {
	case <-time.After(5 * time.Second):
	}

}
