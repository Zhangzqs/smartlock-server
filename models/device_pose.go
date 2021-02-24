package models

import "github.com/beego/beego/v2/client/orm"

//DevicePose 设备姿态表
type DevicePose struct {
	DeviceID string `json:"device_id" orm:"pk;unique;column(device_id)"` //设备ID

	//设备的欧拉角姿态
	Row   float32 `json:"row"`
	Pitch float32 `json:"pitch"`
	Yaw   float32 `json:"yaw"`

	//设备温度
	Temperature float32 `json:"temperature"`
}

//NormalDevicePose 正常设备的欧拉角姿态
type NormalDevicePose struct {
	DeviceID string `json:"device_id" orm:"pk;unique;column(device_id)"` //设备ID

	Row   float32 `json:"row"`
	Pitch float32 `json:"pitch"`
	Yaw   float32 `json:"yaw"`
}

func init() {
	orm.RegisterModel(
		new(DevicePose),
		new(NormalDevicePose),
	)
}
