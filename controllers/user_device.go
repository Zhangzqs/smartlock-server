package controllers

import (
	"encoding/json"
	"smartlock-server/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// /user/device 控制器
type UserDevicesController struct {
	beego.Controller
}

//获取该用户绑定的所有设备ID
func (c *UserDevicesController) Get() {
	var user models.User
	user.UserName = c.Ctx.Input.Query("user_name")
	o := orm.NewOrm()
	if o.Read(&user) != nil {
		logs.Warn("用户信息读取出错", user.UserName)
		c.Ctx.WriteString("用户信息读取出错")
		return
	}
	//此时已经正确读取

	userName := user.UserName
	qs := o.QueryTable("user_device")
	var result []models.UserDevice
	_, _ = qs.Filter("user_name", userName).All(&result)

	c.Data["json"] = result
	_ = c.ServeJSON()
	logs.Debug(result)
}

//绑定一个新的设备
func (c *UserDevicesController) Post() {
	var userDevice models.UserDevice
	userDevice.Auth = models.AdminAuth

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userDevice)
	if err != nil {
		logs.Error("JSON文件解析失败")
		logs.Error(err)
		c.Ctx.WriteString("JSON文件解析失败")
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable("user_device")
	if qs.Filter("user_name", userDevice.UserName).Filter("device_id", userDevice.DeviceID).Exist() {
		//该设备已完成绑定
		c.Ctx.WriteString("该设备已绑定")
		return
	}

	//未绑定
	_, _ = o.Insert(&userDevice)
	c.Ctx.WriteString("成功绑定")
}

//解绑一个设备
func (c *UserDevicesController) Delete() {
	var userDevice models.UserDevice
	userDevice.Auth = models.AdminAuth

	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &userDevice)

	o := orm.NewOrm()
	qs := o.QueryTable("user_device")

	a := qs.Filter("user_name", userDevice.UserName).Filter("device_id", userDevice.DeviceID)

	if a.Exist() {
		_ = a.One(&userDevice)
		return
	}
	_, _ = o.Delete(&userDevice)
}
