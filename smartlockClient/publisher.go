package smartlockClient

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 控制舵机
func controlServo(client mqtt.Client, deviceId string,agree int)  {
	topic := fmt.Sprintf("/smartlock/lock/%s/servo", deviceId)
	payload, _ := json.Marshal(
		struct {
			Agree int `json:"agree"`
		}{Agree: agree},
	)
	token := client.Publish(topic, 2, false, payload)
	logs.Debug("发布消息，","话题：",topic, string(payload))
	token.Wait()
}

// 开锁
func Unlock( deviceId string) {
	logs.Debug("尝试发布消息: 开锁")
	controlServo(client,deviceId,-1)
}

// 锁定门锁(门锁常开)
func LockOpen(deviceId string)  {
	logs.Debug("尝试发布消息：门锁常开")
	controlServo(client,deviceId,-3)
}

// 远程重启
func Restart(deviceId string) {
	logs.Debug("尝试发布消息: 重启")
	topic := fmt.Sprintf("/smartlock/lock/%s/esp", deviceId)
	payload, _ := json.Marshal(
		struct {
			Msg string `json:"msg"`
		}{Msg: "restart"},
	)
	token := client.Publish(topic, 2, false, payload)
	logs.Debug("Publish message", string(payload))
	token.Wait()
}
