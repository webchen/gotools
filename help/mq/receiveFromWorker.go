package mq

import (
	"fmt"
	"time"

	"github.com/webchen/gotools/base"

	"github.com/webchen/gotools/base/conf"
	"github.com/webchen/gotools/help/logs"
	"github.com/webchen/gotools/help/tool/nettool"

	"github.com/streadway/amqp"
)

// Receiver  定义接收者接口
type Receiver interface {
	Consumer([]byte) error
}

var connStr, readExchangeName, readExchangeType, readQueueName, readKey string

func init() {
	connStr = (conf.GetConfig("rabbitmq.connectionString", "")).(string)
	if connStr == "" {
		logs.Critial("rabbitmq 配置为空")
		return
	}

	readExchangeName = (conf.GetConfig("rabbitmq.read.exchange", "data_to_gateway")).(string)
	readExchangeType = (conf.GetConfig("rabbitmq.read.exchangeType", "fanout")).(string)
	readQueueName = (conf.GetConfig("rabbitmq.read.queue", "data")).(string)
	readKey = nettool.GetLocalFirstIPStr()
	// 如果是direct，直接走readKey，队列名必须置为空
	if readExchangeType == "direct" {
		readQueueName = ""
	}
}

func restart(receiver Receiver) {
	time.Sleep(time.Second * 2)
	logs.Info("restart receiver ...")
	receiveData(receiver)
}

// ReceiveWorkerData 接收数据
func ReceiveWorkerData(receiver Receiver) {
	num := int(conf.GetConfig("rabbitmq.read.num", 2.0).(float64))
	for j := 0; j < num; j++ {
		receiveData(receiver)
	}
}

func receiveData(receiver Receiver) {
	base.Go(func() {
		defer func() {
			if p := recover(); p != nil {
				logs.Error("获取MQ对象的协程挂掉了", p)
				restart(receiver)
				return
			}
		}()
		connConsumer, err1 := amqp.Dial(connStr)
		if logs.ErrorProcess(err1, "rabbitmq connect error") {
			restart(receiver)
			return
		}

		channelConsumer, err1 := connConsumer.Channel()
		if logs.ErrorProcess(err1, "consumer channel error") {
			restart(receiver)
			return
		}
		defer (func() {
			channelConsumer.Close()
			connConsumer.Close()
		})()
		if connConsumer.IsClosed() {
			logs.Error("receive connection was closed ...", nil)
			restart(receiver)
			return
		}
		q, err := channelConsumer.QueueDeclare(
			readQueueName, // name
			true,          // durable
			false,         // delete when unused
			false,         // exclusive
			true,          // no-wait
			nil,           // arguments
		)
		if logs.ErrorProcess(err, "读取MQ的数据失败") {
			restart(receiver)
			return
		}
		channelConsumer.ExchangeDeclare(readExchangeName, readExchangeType, true, false, false, true, nil)
		channelConsumer.QueueBind(q.Name, readKey, readExchangeName, true, nil)
		channelConsumer.Qos(5, 0, true) // 1分话，确保rabbitmq会一个一个发消息。5的话，表示预读5个
		closeChan := make(chan *amqp.Error, 1)
		notifyClose := channelConsumer.NotifyClose(closeChan)
		msgs, err := channelConsumer.Consume(
			q.Name,    // queue
			"gateway", // consumer
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			true,      // no-wait
			nil,       // args
		)
		if logs.ErrorProcess(err, fmt.Sprintf("读取MQ数据失败，channel关闭  [%+v]", q)) {
			restart(receiver)
			return
		}
		for {
			select {
			case e := <-notifyClose:
				logs.Error("channel error", e.Error())
				restart(receiver)
				return
			case d := <-msgs:
				err := receiver.Consumer(d.Body)
				logs.ErrorProcess(err, fmt.Sprintf("消费队列出错 [%s]", d.Body))
				logs.Message(fmt.Sprintf("after processWorker   MessageID : [%s]", d.MessageId), string(d.Body), false)
			}
		}
	})
}
