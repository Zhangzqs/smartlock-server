package models

import "github.com/beego/beego/v2/client/orm"

// 设备姿态表
type DevicePose struct {
	DeviceId string		`json:"device_id" orm:"pk;unique"`	//设备ID
	
	//设备的欧拉角姿态
	Row	 	float32		`json:"row"`
	Pitch 	float32		`json:"pitch"`
	Yaw		float32		`json:"yaw"`

	//设备温度
	Temperature	float32	`json:"temperature"`
}

func init() {
	orm.RegisterModel(new(DevicePose))
}