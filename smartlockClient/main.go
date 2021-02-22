package smartlockClient

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client

// 连接到mqtt服务器
func ConnectToServer() mqtt.Client {
	broker := "broker.emqx.io"
	port := 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d",broker,port))
	opts.SetClientID("SmartLock_Server")
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		logs.Debug("Received unhandled message: ",string(message.Payload())," from topic",message.Topic())
	})
	opts.SetAutoReconnect(true);
	opts.OnConnect = func(client mqtt.Client) {
		logs.Debug("Mqtt server connected.")
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logs.Debug("Connect lost:",err)
		ConnectToServer()
	}

	client := mqtt.NewClient(opts)

	if token:=client.Connect();token.Wait()&&token.Error()!=nil{
		panic(token.Error())
	}

	return client
}



func init() {
	client = ConnectToServer()
	RegisterForMqttRouter(client)
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
