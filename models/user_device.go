package models

import "github.com/beego/beego/v2/client/orm"

const (
	// AdminAuth 管理员字段常量
	AdminAuth = "admin"
)

//UserDevice 用户和设备绑定的关系表
type UserDevice struct {
	ID       int    `json:"id" orm:"column(id)"`               //自增索引，默认作为主键
	UserName string `json:"user_name"`                         //用户名
	DeviceID string `json:"device_id" orm:"column(device_id)"` //设备id
	Auth     string `json:"admin"`                             //是否为管理员
	Info     string `json:"info"`                              //备注信息
}

func init() {
	orm.RegisterModel(new(UserDevice))
}
