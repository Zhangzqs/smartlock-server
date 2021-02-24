package smartlockClient

import (
	"encoding/json"
	"fmt"
	"smartlock-server/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//onUidReceived 当收到一张卡片的UID消息时
func onUidReceived(client mqtt.Client, message mqtt.Message) {
	logs.Debug("收到mqtt消息", message)
	var uidMessage struct {
		DeviceID string `json:"device_id"`
		UID      string `json:"uid"`
	}

	_ = json.Unmarshal(message.Payload(), &uidMessage)

	//先根据uid找到其绑定的用户名
	var cardUser models.CardUser
	cardUser.UID = uidMessage.UID
	o := orm.NewOrm()
	err := o.Read(&cardUser)
	if err != nil {
		logs.Warn("开锁失败，用户信息读取出错,不存在该UID的用户", cardUser.UID)
		UserUnlockLog(
			uidMessage.DeviceID,
			"",
			models.CardMethod,
			false,
			fmt.Sprintf("门卡UID:%s", cardUser.UID),
		)
		return
	}

	//读取成功

	var userDevice models.UserDevice
	qs := o.QueryTable("user_device")
	err = qs.Filter("user_name", cardUser.UserName).
		Filter("device_id", uidMessage.DeviceID).One(&userDevice)

	if err == orm.ErrNoRows {
		//找不到设备
		logs.Warn("设备:", uidMessage.DeviceID, "用户：", cardUser.UserName, "不存在关系，无权限")
		logs.Warn(err)
		UserUnlockLog(uidMessage.DeviceID, cardUser.UserName, models.CardMethod, false, "用户与设备无关联")
		return
	}

	//有开锁权限，那么就下发MQTT指令开锁并记录开锁日志
	Unlock(uidMessage.DeviceID)

	UserUnlockLog(uidMessage.DeviceID, cardUser.UserName, models.CardMethod, true, "正常开锁")

	logs.Debug("成功开锁")
}

// 当接收到一个按键请求时
func onButtonReceived(client mqtt.Client, message mqtt.Message) {
	logs.Debug("接收到一个按键请求", string(message.Payload()))
	var buttonEvent struct {
		DeviceID string `json:"device_id"`
		Type     string `json:"type"`
	}

	_ = json.Unmarshal(message.Payload(), &buttonEvent)

	switch buttonEvent.Type {
	case "click":
		// 点击按钮可开锁
		ButtonUnlockLog(buttonEvent.DeviceID, models.UnlockType)
	case "double_click":
		// 双击按钮可门锁常开
		ButtonUnlockLog(buttonEvent.DeviceID, models.OpenType)
	case "long_press":
		// 长按按钮可重启设备
		logs.Debug("设备", buttonEvent.DeviceID, "正常重启")
	}
}

// 当接收到某个设备的在线/离线状态时
func onStatusReceived(client mqtt.Client, message mqtt.Message) {
	var deviceStatus models.DeviceStatus
	err := json.Unmarshal(message.Payload(), &deviceStatus)
	if err != nil {
		logs.Warn("收到一个非法json格式的deviceStatus", err)
		return
	}

	o := orm.NewOrm()
	_, _ = o.Raw(
		"REPLACE INTO `device_status` VALUES(?,?)",
		deviceStatus.DeviceID,
		deviceStatus.Status,
	).Exec()
	logs.Debug("设备在线状态：", deviceStatus)
}

// 当接收到某个设备的姿态信息时
func onPoseReceived(client mqtt.Client, message mqtt.Message) {
	var devicePose models.DevicePose
	err := json.Unmarshal(message.Payload(), &devicePose)
	if err != nil {
		logs.Warn("收到一个非法json格式的devicePose", err)
		return
	}

	o := orm.NewOrm()
	_, _ = o.Raw(
		"REPLACE INTO `device_pose` VALUES(?,?,?,?,?)",
		devicePose.DeviceID,
		devicePose.Row,
		devicePose.Pitch,
		devicePose.Yaw,
		devicePose.Temperature,
	).Exec()
	logs.Debug("设备姿态：", devicePose)
}

// 注册所有的MQTT话题路由
func RegisterForMqttRouter(client mqtt.Client) {
	prefix := "/smartlock/server"
	client.Subscribe(prefix+"/uid", 2, onUidReceived)
	client.Subscribe(prefix+"/button", 2, onButtonReceived)
	client.Subscribe(prefix+"/device_status", 2, onStatusReceived)
	client.Subscribe(prefix+"/device_pose", 2, onPoseReceived)
}
