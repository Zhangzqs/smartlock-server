package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"smartlock-server/locklog"
	"smartlock-server/models"
	"smartlock-server/smartlockClient"
	"time"
)

// /token_unlock 控制器
type TokenUnlockController struct {
	beego.Controller
}

func (c *TokenUnlockController) Post() {
	var tokenMsg struct {
		DeviceID string `json:"device_id"`
		Token    string `json:"token"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &tokenMsg)

	//按设备，按卡号过滤
	o := orm.NewOrm()
	qs := o.QueryTable("token")
	resultView := qs.Filter("device_id", tokenMsg.DeviceID).
		Filter("token", tokenMsg.Token)

	if !resultView.Exist() {
		//该门锁上不存在该门卡
		logs.Warn("开锁失败，不存在该口令", tokenMsg.Token)
		locklog.UserUnlockLog(
			tokenMsg.DeviceID,
			"",
			models.TokenMethod,
			false,
			tokenMsg.Token,
			"不存在的口令",
		)
		return
	}

	var tokenModel models.Token
	_ = resultView.One(&tokenModel)

	// 存在该门卡，开始查看有效期
	//判断门卡有效期
	if now := int(time.Now().Unix()); now < tokenModel.BeginTime {
		logs.Warn("开锁失败，门卡未生效")
		locklog.UserUnlockLog(tokenModel.DeviceID, tokenModel.UserName, models.CardMethod, false, tokenMsg.Token, "未生效的口令")
		return
	} else {
		//开始检测门卡是否失效
		if tokenModel.EndTime != 0 && now > tokenModel.EndTime {
			//门卡已失效
			logs.Warn("开锁失败，门卡已过期")
			locklog.UserUnlockLog(tokenModel.DeviceID, tokenModel.UserName, models.CardMethod, false, tokenMsg.Token, "已失效的口令")
			return
		}
	}

	//有开锁权限，那么就下发MQTT指令开锁并记录开锁日志
	smartlockClient.Unlock(tokenModel.DeviceID)

	locklog.UserUnlockLog(tokenModel.DeviceID, tokenModel.UserName, models.CardMethod, true, tokenMsg.Token, "正常开锁")

	logs.Debug("成功开锁")
	c.Ctx.WriteString("OK")
}
