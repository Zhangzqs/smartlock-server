package routers

import (
	"smartlock-server/controllers"
	beego "github.com/beego/beego/v2/server/web"
)



func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/user",&controllers.UsersController{})
    beego.Router("/user/:user_name",&controllers.UserController{})
    beego.Router("/user/device",&controllers.UserDevicesController{})
    beego.Router("/user/card",&controllers.CardController{})
    beego.Router("/device/:device_id",&controllers.DeviceController{})
	beego.Router("/login",&controllers.LoginController{})
    beego.Router("/register",&controllers.RegisterController{})
}
