package mq

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	base "github.com/webchen/gotools/base"
	"github.com/webchen/gotools/base/conf"
	"github.com/webchen/gotools/base/jsontool"
	"github.com/webchen/gotools/help/logs"
	"github.com/webchen/gotools/help/util/goqueue"

	"github.com/streadway/amqp"
)

var connWorkerStr, sendWorkerExchangeName, sendWorkerExchangeType, sendWorkerQueueName, sendWorkerKey string

var connectNum, channelNum int

var connBaseChannelNum = 100

// connection 列表
var connSendList sync.Map

// channel 列表
//var channelSendList sync.Map
var channelSendList = goqueue.NewQueue()

var reConnectLock = base.NewTryMutex()

// SendFormat 发送数据的格式
type SendFormat struct {
	//	Timestamp int64  `json:"timestamp"`
	T     uint8  `json:"t"`
	Trace string `json:"trace"`
	Data  string `json:"data"`
}

// ChannelObject channel对象
type ChannelObject struct {
	ChannelObj    *amqp.Channel
	ChannelID     int
	ConnectionObj *amqp.Connection
	ConnectionID  int
}

func init() {
	connWorkerStr = (conf.GetConfig("rabbitmq.connectionString", "")).(string)
	if connWorkerStr == "" {
		logs.Critial("rabbitmq 配置为空")
		return
	}
	sendWorkerExchangeName = (conf.GetConfig("rabbitmq.send.exchange", "data_to_worker")).(string)
	sendWorkerExchangeType = (conf.GetConfig("rabbitmq.send.exchangeType", "fanout")).(string)
	sendWorkerQueueName = (conf.GetConfig("rabbitmq.send.queue", "data")).(string)
	sendWorkerKey = ""
	connectNum = int(conf.GetConfig("rabbitmq.send.connectNum", 4.0).(float64))
	channelNum = int(conf.GetConfig("rabbitmq.send.channelNum", 10.0).(float64))
	if base.IsBuild() {
		return
	}
	initPool(1, connectNum)
}

func reConnect() {
	fmt.Println("reConnected mq ..")
	if reConnectLock.IsLocked() {
		return
	}
	reConnectLock.Lock()
	initPool(connectNum+1, connectNum*2)
}

func initConnection(start int, max int) {
	for j := start; j <= max; j++ {
		conn, err := amqp.Dial(connWorkerStr)
		if logs.ErrorProcess(err, "无法创建MQ连接") {
			continue
		}
		conn.Properties["ConnectionName"] = sendWorkerQueueName
		conn.Properties["connection_name"] = sendWorkerQueueName
		connSendList.Store(j, conn)
	}
}

func initChannel(connStart int, connMax int) {
	connSendList.Range(func(key, value interface{}) bool {

		if key.(int) < connStart || key.(int) > connMax {
			return true
		}

		conn := value.(*amqp.Connection)
		if conn.IsClosed() {
			return true
		}

		for i := 1; i <= channelNum; i++ {
			c, err := conn.Channel()
			if logs.ErrorProcess(err, "无法创建channel") {
				continue
			}
			obj := &ChannelObject{}
			obj.ChannelObj = c
			obj.ChannelID = key.(int)*connBaseChannelNum + i
			obj.ConnectionObj = conn
			obj.ConnectionID = key.(int)
			channelSendList.Push(obj)
		}
		return true
	})
}

func initPool(start int, max int) {
	initConnection(start, max)
	initChannel(start, max)
}

func closeAll() {
	list := channelSendList.Clear2List()
	for _, v := range list {
		val := v.(*ChannelObject)
		val.ChannelObj.Close()
		if !val.ConnectionObj.IsClosed() {
			val.ConnectionObj.Close()
		}
	}
	connSendList.Range(func(k, v interface{}) bool {
		if !v.(*amqp.Connection).IsClosed() {
			v.(*amqp.Connection).Close()
		}
		connSendList.Delete(k)
		return true
	})
}

func getChannle() *ChannelObject {
	//logs.Warning("mq queue before pop ", channelSendList.Len(), false)
	list, err := channelSendList.Pop()
	if logs.ErrorProcess(err, "无法获取channel对象") {
		return nil
	}
	obj, ok := list.(*ChannelObject)
	if !ok {
		logs.Error("channel对象转换失败", list)
		return nil
	}

	if obj.ConnectionObj.IsClosed() {
		obj.ConnectionObj, err = amqp.Dial(connWorkerStr)
		if logs.ErrorProcess(err, "重连mq的connection") {
			return nil
		}
	}
	/*
		if obj.ChannelObj.IsClosed() {
			obj.ChannelObj, err = obj.ConnectionObj.Channel()
			if logs.ErrorProcess(err, "重连mq的channel") {
				return nil
			}
		}
	*/
	return obj
}

// 实际发送的函数，解耦内容和其他附加参数（重试等）
func send(data *SendFormat) {
	if data.T > 3 {
		return
	}

	randm := strconv.FormatInt(time.Now().UnixNano(), 10)
	var err error
	obj := getChannle()
	if obj == nil {
		time.Sleep(time.Millisecond * 1000)
		reConnect()
		send(data)
		return
	}
	defer func() {
		channelSendList.Push(obj)
	}()
	ch := obj.ChannelObj
	_, err = ch.QueueDeclare(
		sendWorkerQueueName, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		true,                // no-wait
		nil,                 // arguments
	)
	//	fmt.Printf("get mq obj [%+v]\n", ch)
	if logs.ErrorProcess(err, fmt.Sprintf("Failed to declare a queue. MessageID : [%s]", randm)) {
		time.Sleep(time.Millisecond * 200)
		send(data)
		return
	}
	body := jsontool.MarshalToString(data.Data)
	err = ch.QueueBind(sendWorkerQueueName, sendWorkerKey, sendWorkerExchangeName, true, nil)
	if logs.ErrorProcess(err, fmt.Sprintf("绑定Queue失败， MessageID : [%s]", randm)) {
		time.Sleep(time.Millisecond * 200)
		send(data)
		return
	}
	//	fmt.Printf("bind mq err: [%+v]\n", err)
	err = ch.Publish(
		sendWorkerExchangeName, // exchange
		sendWorkerKey,          // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
			MessageId:    randm,
			Timestamp:    time.Now(),
		})

	if logs.ErrorProcess(err, "Failed Mq Send") {
		if data.T > 3 {
			logs.Warning("该消息重发3次失败", data, false)
			return
		}
		time.Sleep(time.Second)
		data.T++
		base.Go(send, data)
		return
	}
	logs.Message(fmt.Sprintf("sendDataToWorker MessageID : [%s]", randm), data, false)

}

// SendData2Worker 发消息给MQ
func SendData2Worker(data string) {
	if strings.TrimSpace(data) == "" {
		return
	}
	d := &SendFormat{}
	d.T = 1
	d.Data = data
	send(d)
}
