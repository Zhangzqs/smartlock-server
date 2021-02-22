package models

import "github.com/beego/beego/v2/client/orm"

const (
	Offline = "online"
	Online = "offline"
)

// 设备状态表
type DeviceStatus struct {
	DeviceId string		`json:"device_id" orm:"pk;unique"`	//设备ID
	Status	 string		`json:"status"`						//在线/离线的状态
}

func init() {
	orm.RegisterModel(new(DeviceStatus))
}