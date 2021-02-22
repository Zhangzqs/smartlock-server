package controllers

import beego "github.com/beego/beego/v2/server/web"

// /user 控制器
type UsersController struct {
	beego.Controller
}

//获取所有用户
func (c *UsersController) Get() {
	c.Ctx.WriteString("暂不支持获取所有用户")
}
