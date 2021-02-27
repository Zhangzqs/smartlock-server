package locklog

import (
	"smartlock-server/models"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// 记录按键开锁日志
func ButtonUnlockLog(deviceId string, _type string) {
	var log models.ButtonUnlockLog
	log.DeviceID = deviceId
	log.Time = int(time.Now().Unix())
	log.Type = _type
	o := orm.NewOrm()
	_, err := o.Insert(&log)
	if err != nil {
		logs.Warn("记录按键开锁日志发生异常", err)
		return
	}
	logs.Debug("按键开锁日志记录成功", "DeviceID:", deviceId)
}

// UserUnlockLog 记录用户开锁日志
func UserUnlockLog(deviceID string, userName string, method string, success bool, auth string, info string) {
	var log models.UserUnlockLog
	log.UserName = userName
	log.DeviceID = deviceID
	log.Method = method
	log.Info = info
	log.Auth = auth
	if success {
		log.Success = 1
	} else {
		log.Success = 0
	}
	log.Time = int(time.Now().Unix())
	o := orm.NewOrm()
	_, err := o.Insert(&log)
	if err != nil {
		logs.Warn("记录用户开锁日志发生异常", err)
		return
	}
	logs.Debug("用户开锁日志记录成功", "DeviceID:", deviceID, "UserName:", userName, "Method:", method)

}
