package models

import "github.com/beego/beego/v2/client/orm"

const (
	//Offline 离线
	Offline = "offline"

	//Online 在线
	Online = "online"
)

//DeviceStatus 设备状态表
type DeviceStatus struct {
	DeviceID string `json:"device_id" orm:"pk;unique;column(device_id)"` //设备ID
	Status   string `json:"status"`                                      //在线/离线的状态
}

func init() {
	orm.RegisterModel(new(DeviceStatus))
}
