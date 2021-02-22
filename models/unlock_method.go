package models

import "github.com/beego/beego/v2/client/orm"

type base struct {
	UserName string `json:"user_name"`           //开放该开锁方式的用户
	Info     string `json:"info"`                //备注
	BeginTime	 int	`json:"begin_time"`		//有效期开始时间
	EndTime	 int	`json:"end_time"`			//有效期截止时间
}

//卡片开锁
type CardUser struct {
	base
	Uid      string `json:"uid" orm:"pk;unique"` //ID信息
}

//口令开锁
type TokenUser struct {
	base
	Token    string `json:"uid" orm:"pk;unique"` //令牌信息
}

//人脸开锁
type FaceUser struct {
	base
	Token 	 string	`json:"token" orm:"pk;unique"`
}

func init() {
	orm.RegisterModel(
		new(CardUser),
		new(TokenUser),
		new(FaceUser),
	)
}
