package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"smartlock-server/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/user", &controllers.UsersController{})
	beego.Router("/user/:user_name", &controllers.UserController{})
	beego.Router("/device", &controllers.UserDevicesController{})
	beego.Router("/card", &controllers.CardController{})
	beego.Router("/token", &controllers.TokenController{})
	beego.Router("token_unlock", &controllers.TokenUnlockController{})
	beego.Router("/device/:device_id", &controllers.DeviceController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})
}
