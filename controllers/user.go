package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"smartlock-server/models"
)

//	/user/{username} 控制器
type UserController struct {
	beego.Controller
}

// 获取用户信息
func (c *UserController) Get()  {
	var user models.User
	user.UserName = c.Ctx.Input.Param(":user_name")
	o := orm.NewOrm()
	err := o.Read(&user)
	if err != nil {
		logs.Warn("用户信息读取出错",err)
		c.Ctx.WriteString("用户信息读取出错")
		return
	}

	//用户信息成功读取,以json返回
	c.Data["json"] = &user
	_ = c.ServeJSON()
}

//修改用户信息
func (c *UserController) Put()  {
	var user models.User
	user.UserName = c.Ctx.Input.Param(":user_name")
	o := orm.NewOrm()
	if o.Read(&user) != nil {
		logs.Warn("用户信息读取出错")
		c.Ctx.WriteString("用户信息读取出错")
		return
	}
	//此时已经正确读取

	userName := user.UserName

	//现在开始解析json
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &user)
	if err != nil {
		logs.Error("JSON文件解析失败", string(data[:]))
		logs.Error(err)
		responseData := loginRegResponse{Code: JsonFormatError} //1: 表示JSON文件有误
		c.Data["json"] = &responseData
		_ = c.ServeJSON()
		return
	}

	logs.Debug("JSON文件解析成功")

	if userName != user.UserName{
		c.Ctx.WriteString("不支持修改用户名")
		return
	}

	//可以更新了
	_, err = o.Update(&user)
	if err != nil {
		logs.Warn("用户信息更新出错")
		c.Ctx.WriteString("用户信息更新出错")
		return
	}

	//正确更新

	c.Data["json"] = &struct {
		Code int	`json:"code"`
	}{0}
	_ = c.ServeJSON()
}

//删除用户
func (c *UserController) Delete()  {
	logs.Warn("暂不支持删除用户")
	c.Ctx.WriteString("暂不支持删除用户")
}