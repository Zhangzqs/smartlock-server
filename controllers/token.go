package controllers

import (
	"encoding/json"
	"smartlock-server/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// /token 控制器
type TokenController struct {
	beego.Controller
}

//获取该用户拥有的所有卡片
func (c *TokenController) Get() {
	userName := c.Ctx.Input.Query("user_name")
	deviceId := c.Ctx.Input.Query("device_id")

	o := orm.NewOrm()
	qs := o.QueryTable("token")

	var tokens []models.Token

	resultView := qs.Filter("user_name", userName)

	if deviceId != "" {
		resultView = resultView.Filter("device_id", deviceId)
	}

	_, _ = resultView.All(&tokens)

	c.Data["json"] = tokens

	_ = c.ServeJSON()
}

//添加一张新卡片
func (c *TokenController) Post() {
	var token models.Token
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	if err != nil {
		logs.Error("JSON文件解析失败")
		logs.Error(err)
		c.Ctx.WriteString("JSON文件解析失败")
		return
	}
	o := orm.NewOrm()
	_, err = o.Insert(&token)

	logs.Debug("准备添加口令", "用户名:", token.UserName, "口令：", token.Token)
	if err != nil {
		logs.Debug("添加失败")
		c.Ctx.WriteString("失败")
	}
	logs.Debug("添加成功")
	c.Ctx.WriteString("成功")
}

//解绑一张卡片
func (c *TokenController) Delete() {

}
