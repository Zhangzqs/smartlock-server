package main

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/mattn/go-sqlite3"
	_ "smartlock-server/routers"
	_ "smartlock-server/smartlockClient"
)

func main() {
	beego.Run()
}

func init() {
	//数据库配置及初始化

	orm.Debug = false

	//注册数据库
	_ = orm.RegisterDataBase("default", "sqlite3", "database/data.db")

	//自动建表，force如果为true，将会强制更新表，可能会导致数据丢失
	_ = orm.RunSyncdb("default", false, true)

	//日志模块的配置及初始化

	//配置日志输出到控制台
	_ = logs.SetLogger("console")
}
