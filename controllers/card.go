package controllers

import (
	"encoding/json"
	"smartlock-server/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// /card 控制器
type CardController struct {
	beego.Controller
}

//获取该用户拥有的所有卡片
func (c *CardController) Get() {
	userName := c.Ctx.Input.Query("user_name")
	deviceId := c.Ctx.Input.Query("device_id")

	o := orm.NewOrm()
	qs := o.QueryTable("card")

	var cards []models.Card

	resultView := qs.Filter("user_name", userName)

	if deviceId != "" {
		resultView = resultView.Filter("device_id", deviceId)
	}

	_, _ = resultView.All(&cards)

	c.Data["json"] = cards

	_ = c.ServeJSON()
}

//添加一张新卡片
func (c *CardController) Post() {
	var card models.Card
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &card)
	if err != nil {
		logs.Error("JSON文件解析失败")
		logs.Error(err)
		c.Ctx.WriteString("JSON文件解析失败")
		return
	}
	o := orm.NewOrm()
	_, err = o.Insert(&card)

	logs.Debug("准备添加门卡", "用户名:", card.UserName, "卡UID：", card.UID)
	if err != nil {
		logs.Debug("添加失败")
	}
	logs.Debug("添加成功")
	c.Ctx.WriteString("添加成功")
}

//解绑一张卡片
func (c *CardController) Delete() {

}
