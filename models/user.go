package models

import "github.com/beego/beego/v2/client/orm"

type User struct {
	UserName string	`json:"userName" orm:"pk;unique"`	//unique唯一  pk主键
	Password string	`json:"password"`
	Phone string	`json:"phone"`
	QqToken string	`json:"qqToken"`
}

func init() {
	orm.RegisterModel(new(User))
}