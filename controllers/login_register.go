package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"smartlock-server/models"
)

// /login 控制器
type LoginController struct {
	beego.Controller
}

//登陆反馈的结构体
type loginRegResponse struct {
	Code     int    `json:"code"`
	UserName string `json:"user_name"`
}

const (
	Ok              = 0
	JsonFormatError = 1
	NotFoundUser    = 2
	Unknown         = 3
	FoundUser       = 4
	PasswordError   = 5
)

// 登陆控制器逻辑
func (c *LoginController) Post() {
	var user struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	data := c.Ctx.Input.RequestBody

	logs.Debug("收到一个登录请求体：",string(data))

	var responseData loginRegResponse
	err := json.Unmarshal(data, &user)
	if err != nil {
		logs.Error("JSON文件解析失败", string(data[:]))
		logs.Error(err)
		responseData = loginRegResponse{Code: JsonFormatError} //1: 表示JSON文件有误
		c.Data["json"] = &responseData
		_ = c.ServeJSON()
		return
	}

	logs.Debug("JSON文件解析成功")

	o := orm.NewOrm()
	var userModel models.User
	userModel.UserName = user.UserName
	err = o.Read(&userModel)

	logs.Debug("读取用户数据结果: ",err)

	if err != nil {
		if err == orm.ErrNoRows {
			responseData = loginRegResponse{Code: NotFoundUser} //2: 不存在用户
		} else {
			responseData= loginRegResponse{Code: Unknown} //3: 未知原因
			logs.Warn(err)
		}

		c.Data["json"] = &responseData
		_ = c.ServeJSON()
		return
	}

	if userModel.Password != user.Password{
		c.Data["json"] = &loginRegResponse{Code: PasswordError}
		_ = c.ServeJSON()
		return
	}

	//读取成功
	responseData = loginRegResponse{
		Code:     Ok,            //0: 表示成功登录
		UserName: user.UserName, //返回登陆成功的用户名
	}

	logs.Debug("登录结果：",responseData)
	c.Data["json"] = &responseData
	_ = c.ServeJSON()
}

// /register 控制器
type RegisterController struct {
	beego.Controller
}

//注册控制器逻辑
func (c *RegisterController) Post() {
	var user struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	data := c.Ctx.Input.RequestBody

	logs.Debug("收到一个注册请求体：",string(data))

	var responseData loginRegResponse
	err := json.Unmarshal(data, &user)
	if err != nil {
		logs.Error("JSON文件解析失败", string(data[:]))
		logs.Error(err)
		responseData = loginRegResponse{Code: JsonFormatError} //1: 表示JSON文件有误
		c.Data["json"] = &responseData
		_ = c.ServeJSON()
		return
	}

	logs.Debug("JSON文件解析成功")

	o := orm.NewOrm()
	var userModel models.User
	userModel.UserName = user.UserName
	err = o.Read(&userModel)

	logs.Debug("读取用户数据结果: ",err)

	if err == nil{
		//读取无误，说明数据库已存在该用户了
		responseData = loginRegResponse{Code: FoundUser} //4: 已存在该用户
		c.Data["json"] = &responseData
		_ = c.ServeJSON()
		return
	}

	//此时读取是有误的
	if err != orm.ErrNoRows {
		//如果错误原因不是因为不存在用户引起的，那就输出未知错误
		responseData = loginRegResponse{Code: Unknown}
		c.Data["json"] = &responseData
		_ = c.ServeJSON()
		return
	}


	//此时证明该用户不存在，可以开始创建用户

	userModel.Password = user.Password

	id,err := o.Insert(&userModel)

	if err != nil {
		//插入出错
		logs.Warn("创建用户出错",user.UserName,err)
		responseData = loginRegResponse{Code: Unknown}
		c.Data["json"] = &responseData
		_ = c.ServeJSON()
		return
	}

	//成功创建用户
	logs.Debug("成功创建用户",user.UserName,"当前用户数目为：",id)

	responseData = loginRegResponse{
		Code:     Ok,            //0: 表示成功
		UserName: user.UserName, //返回注册成功的用户名
	}

	logs.Debug("注册结果：",responseData)
	c.Data["json"] = &responseData
	_ = c.ServeJSON()

}
