package smartlockClient

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"smartlock-server/locklog"
	"smartlock-server/models"
	"time"
)

//onUidReceived 当收到一张卡片的UID消息时
func onUidReceived(client mqtt.Client, message mqtt.Message) {
	logs.Debug("收到门卡请求", string(message.Payload()))
	var uidMsg struct {
		DeviceID string `json:"device_id"`
		UID      string `json:"uid"`
	}

	_ = json.Unmarshal(message.Payload(), &uidMsg)

	//按设备，按卡号过滤
	o := orm.NewOrm()
	qs := o.QueryTable("card")
	resultView := qs.Filter("device_id", uidMsg.DeviceID).
		Filter("uid", uidMsg.UID)

	if !resultView.Exist() {
		//该门锁上不存在该门卡
		logs.Warn("开锁失败，不存在该门卡", uidMsg.UID)
		locklog.UserUnlockLog(
			uidMsg.DeviceID,
			"",
			models.CardMethod,
			false,
			uidMsg.UID,
			"不存在的门卡",
		)
		return
	}

	var cardModels models.Card
	_ = resultView.One(&cardModels)

	// 存在该门卡，开始查看有效期
	//判断门卡有效期
	if now := int(time.Now().Unix()); now < cardModels.BeginTime {
		logs.Warn("开锁失败，门卡未生效")
		locklog.UserUnlockLog(cardModels.DeviceID, cardModels.UserName, models.CardMethod, false, uidMsg.UID, "未生效的门卡")
		return
	} else {
		//开始检测门卡是否失效
		if cardModels.EndTime != 0 && now > cardModels.EndTime {
			//门卡已失效
			logs.Warn("开锁失败，门卡已过期")
			locklog.UserUnlockLog(cardModels.DeviceID, cardModels.UserName, models.CardMethod, false, uidMsg.UID, "已失效的门卡")
			return
		}
	}

	//有开锁权限，那么就下发MQTT指令开锁并记录开锁日志
	Unlock(cardModels.DeviceID)

	locklog.UserUnlockLog(cardModels.DeviceID, cardModels.UserName, models.CardMethod, true, uidMsg.UID, "正常开锁")

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
		locklog.ButtonUnlockLog(buttonEvent.DeviceID, models.UnlockType)
	case "double_click":
		// 双击按钮可门锁常开
		locklog.ButtonUnlockLog(buttonEvent.DeviceID, models.OpenType)
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
