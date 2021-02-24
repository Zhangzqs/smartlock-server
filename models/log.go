package models

import "github.com/beego/beego/v2/client/orm"

const (
	//CardMethod 门卡开锁
	CardMethod = "card"

	//AppMethod App开锁
	AppMethod = "app" //App开锁

	//FaceMethod 人脸开锁
	FaceMethod = "face" //人脸开锁

	//FingerMethod 指纹解锁
	FingerMethod = "finger" //指纹解锁

	//OpenType 门锁常开
	OpenType = "open"

	//CloseType 门锁常闭
	CloseType = "close"

	//UnlockType 正常开锁
	UnlockType = "unlock"
)

//UserUnlockLog 用户开锁日志表
type UserUnlockLog struct {
	ID       int    `json:"id" orm:"auto;column(id)"`
	UserName string `json:"user_name"`
	DeviceID string `json:"device_id" orm:"column(device_id)"`
	Method   string `json:"method"`
	Success  int    `json:"success"` //success为1表示成功开锁，0表示失败
	Info     string `json:"info"`    //开锁备注
	Time     int    `json:"time"`
}

//ButtonUnlockLog 按键开锁日志表
type ButtonUnlockLog struct {
	ID       int    `json:"id" orm:"auto;column(id)"`
	DeviceID string `json:"device_id" orm:"column(device_id)"`
	Type     string `json:"type"` //开锁类型
	Time     int    `json:"time"`
}

func init() {
	orm.RegisterModel(
		new(UserUnlockLog),
		new(ButtonUnlockLog),
	)
}
