package controllers

import (
	"fmt"
	"smartlock-server/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// /user/card 控制器
type CardController struct {
	beego.Controller
}

//获取该用户拥有的所有卡片
func (c *CardController) Get() {
	var user models.User
	user.UserName = c.Ctx.Input.Query("user_name")
	o := orm.NewOrm()
	if o.Read(&user) != nil {
		logs.Warn("用户信息读取出错")
		c.Ctx.WriteString("用户信息读取出错")
		return
	}
	//此时已经正确读取

	userName := user.UserName
	c.Ctx.WriteString(fmt.Sprintf("读取用户%s的所有卡片", userName))
}

//添加一张新卡片
func (c *CardController) Post() {

}

//解绑一张卡片
func (c *CardController) Delete() {

}
