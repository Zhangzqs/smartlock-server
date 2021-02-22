package models

import "github.com/beego/beego/v2/client/orm"

const (
	AdminAuth = "admin"
)

type UserDevice struct {
	Id 		 int		`json:"id"`			//自增索引，默认作为主键
	UserName string		`json:"user_name"`	//用户名
	DeviceId string		`json:"device_id"`	//设备id
	Auth	 string		`json:"admin"`		//是否为管理员
}

func init() {
	orm.RegisterModel(new(UserDevice))
}