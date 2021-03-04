package models

import "github.com/beego/beego/v2/client/orm"

//CardUser 卡片开锁
type Card struct {
	ID        int    `json:"id" orm:"auto;column(id)"`
	UserName  string `json:"user_name"`                         //开放该开锁方式的用户
	DeviceID  string `json:"device_id" orm:"column(device_id)"` //卡片所生效的设备
	Info      string `json:"info"`                              //备注
	BeginTime int    `json:"begin_time"`                        //有效期开始时间
	EndTime   int    `json:"end_time"`                          //有效期截止时间
	UID       string `json:"uid" orm:"column(uid)"`             //ID信息
}

//TokenUser 口令开锁
type Token struct {
	ID        int    `json:"id" orm:"auto;column(id)"`
	UserName  string `json:"user_name"`                         //开放该开锁方式的用户
	DeviceID  string `json:"device_id" orm:"column(device_id)"` //口令所生效的设备
	Info      string `json:"info"`                              //备注
	BeginTime int    `json:"begin_time"`                        //有效期开始时间
	EndTime   int    `json:"end_time"`                          //有效期截止时间
	Token     string `json:"token"`                             //令牌信息
}

//FaceUser 人脸开锁
// type FaceUser struct {
// 	base
// 	Token string `json:"token" orm:"pk;unique"`
// }

func init() {
	orm.RegisterModel(
		new(Card),
		new(Token),
		//new(FaceUser),
	)
}
