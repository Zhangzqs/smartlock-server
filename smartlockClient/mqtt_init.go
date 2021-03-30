package smartlockClient

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"time"
)

var client mqtt.Client

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// 连接到mqtt服务器
func ConnectToServer() mqtt.Client {
	//broker := "broker.emqx.io"
	broker := "10.1.160.240"
	port := 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("SmartLock_Server_" + RandString(10))
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		logs.Debug("Received unhandled message: ", string(message.Payload()), " from topic", message.Topic())
	})
	opts.OnConnect = func(client mqtt.Client) {
		logs.Debug("Mqtt server connected.")
		//注册mqtt话题
		RegisterForMqttRouter(client)
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logs.Debug("Connect lost:", err)
		client.Disconnect(250)
		logs.Debug("Mqtt disconnected")
		//ConnectToServer()
	}
	//使用SetAutoReconnect代替手动重连
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(2 * time.Second)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		//panic(token.Error())
		logs.Warn(token.Error())
	}
	return client
}

func init() {
	client = ConnectToServer()
	////定时发布话题
	//go func(client mqtt.Client) {
	//	topic:="Server1"
	//	for i:=0;i<9999999;i++{
	//		text:=fmt.Sprintf("Message %d",i)
	//		token:=client.Publish(topic,0,false,text)
	//		logs.Debug("Publish message",text)
	//		token.Wait()
	//		time.Sleep(time.Second)
	//	}
	//}(client)

}
