package controllers

import (
	json "encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"smartlock-server/models"
	"smartlock-server/smartlockClient"
)

//	/device/{device_id} 控制器
type DeviceController struct {
	beego.Controller
}

// 未注册信息，密码键盘，获取设备信息
func (c *DeviceController) Get()  {
	deviceId := c.Ctx.Input.Param(":device_id")
	c.Data["DeviceId"] = deviceId
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.UserDevice))
	exist := qs.Filter("device_id",deviceId).Exist()

	if !exist {	//未注册的设备
		// 使用浏览器访问时
		logs.Warn("检测到一个未注册的设备: ", deviceId)
		c.TplName = "not_found_device.html"
		return
	}


	//这是一个已经记录在案的设备，判断是否存在field查询字段
	field:=c.Ctx.Input.Query("field")
	if field==""{
		//空字符串，那就是网页访问,显示令牌输入界面
		c.TplName = "token_input.html"
		return
	}

	//能执行到这里，说明就是普通的api调用，用于查询设备信息

	//TODO

}


// 向设备发送消息
func (c *DeviceController) Post()  {
	var device models.DeviceStatus
	device.DeviceId = c.Ctx.Input.Param(":device_id")
	o := orm.NewOrm()
	err := o.Read(&device)
	if err != nil {
		logs.Warn("设备信息读取出错",err)
		c.Ctx.WriteString("设备信息读取出错")
		return
	}

	//设备信息成功读取,开始判断是否在线
	if device.Status == models.Offline{
		//如果设备离线了
		logs.Warn("设备",device.DeviceId,"已离线")
		c.Ctx.WriteString("设备离线")
		return
	}

	//如果此时设备还在线，那么就开始处理消息

	var ins struct{
		UserName 		string	`json:"user_name"`	//控制者
		Dist			string	`json:"dist"`		//目标
		Instruction 	string	`json:"instruction"`	//指令
	}

	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &ins)

	switch ins.Dist {
	case "servo":
		switch ins.Instruction {
		case "unlock":
		case "open":
		}
	case "led":
	case "sound":
	case "esp":
		switch ins.Instruction {
		case "restart":
			smartlockClient.Unlock(device.DeviceId)
		}
	}

	
	c.Ctx.WriteString("设备控制成功")
}