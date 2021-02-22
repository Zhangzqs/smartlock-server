package models

import "github.com/beego/beego/v2/client/orm"

const (
	//Method字段常量
	CardMethod = "card"			//门卡开锁
	AppMethod = "app"			//App开锁
	FaceMethod = "face"			//人脸开锁
	FingerMethod = "finger"		//指纹解锁

	//Type字段常量
	OpenType   = "open"			//常开
	CloseType  = "close"		//关闭
	UnlockType = "unlock"		//开锁
)

// 用户开锁日志表
type UserUnlockLog struct {
	Id		 int	 `json:"id" orm:"auto"`
	UserName string `json:"user_name"`
	DeviceId string `json:"device_id"`
	Method   string    `json:"method"`
	Success int		`json:"success"`	//success为1表示成功开锁，0表示失败
	Info 	string	`json:"info"`		//开锁备注
	Time     int    `json:"time"`
}

//按键开锁日志表
type ButtonUnlockLog struct {
	Id 		 int		`json:"id" orm:"auto"`
	DeviceId string		`json:"device_id"`
	Type	 string 	`json:"type"`		//开锁类型
	Time     int		`json:"time"`
}

func init() {
	orm.RegisterModel(
		new(UserUnlockLog),
		new(ButtonUnlockLog),
	)
}
