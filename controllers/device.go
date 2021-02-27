package controllers

import (
	json "encoding/json"
	"smartlock-server/locklog"
	"smartlock-server/models"
	"smartlock-server/smartlockClient"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

//	/device 控制器
type DeviceController struct {
	beego.Controller
}

// 未注册信息，密码键盘，获取设备信息
func (c *DeviceController) Get() {
	deviceID := c.Ctx.Input.Param(":device_id")
	c.Data["DeviceID"] = deviceID
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.UserDevice))
	exist := qs.Filter("device_id", deviceID).Exist()

	if !exist { //未注册的设备
		// 使用浏览器访问时
		logs.Warn("检测到一个未注册的设备: ", deviceID)
		c.TplName = "not_found_device.html"
		return
	}

	//这是一个已经记录在案的设备，显示令牌输入页面
	c.TplName = "token_input.html"
	return
}

// 向设备发送消息
func (c *DeviceController) Post() {
	var device models.DeviceStatus
	device.DeviceID = c.Ctx.Input.Param(":device_id")
	o := orm.NewOrm()
	err := o.Read(&device)
	if err != nil {
		logs.Warn("设备信息读取出错", err)
		c.Ctx.WriteString("设备信息读取出错")
		return
	}

	//设备信息成功读取,开始判断是否在线
	if device.Status == models.Offline {
		//如果设备离线了
		logs.Warn("设备", device.DeviceID, "已离线")
		c.Ctx.WriteString("设备离线")
		return
	}

	//如果此时设备还在线，那么就开始处理消息

	var ins struct {
		UserName    string `json:"user_name"`   //控制者
		Dist        string `json:"dist"`        //目标
		Instruction string `json:"instruction"` //指令
	}

	err = json.Unmarshal(c.Ctx.Input.RequestBody, &ins)

	if err != nil {
		c.Ctx.WriteString("JSON控制指令格式有误")
		return
	}

	switch ins.Dist {

	case "servo":
		switch ins.Instruction {
		case "unlock":
			smartlockClient.Unlock(device.DeviceID)
			locklog.UserUnlockLog(device.DeviceID, ins.UserName, models.AppMethod, true, "正常远程开锁")
		case "open":
			smartlockClient.LockOpen(device.DeviceID)
			locklog.UserUnlockLog(device.DeviceID, ins.UserName, models.AppMethod, true, "正常远程常开门锁")
		}
	case "led":
	case "sound":
	case "esp":
		switch ins.Instruction {
		case "restart":
			smartlockClient.Restart(device.DeviceID)
		case "get_pose":
			var devicePose models.DevicePose
			devicePose.DeviceID = device.DeviceID
			err = o.Read(&devicePose)
			if err != nil {
				c.Ctx.WriteString("无法获取设备的姿态信息")
				return
			}
			c.Data["json"] = &devicePose
			c.ServeJSON()
			return
		}
	}

	c.Ctx.WriteString("设备控制成功")
}
