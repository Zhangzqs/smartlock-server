package smartlockClient

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"smartlock-server/models"
	"time"
)

// 记录按键开锁日志
func ButtonUnlockLog(deviceId string,_type string) {
	var log models.ButtonUnlockLog
	log.DeviceId = deviceId
	log.Time = int(time.Now().Unix())
	log.Type = _type;
	o := orm.NewOrm()
	_, err := o.Insert(&log)
	if err != nil {
		logs.Warn("记录按键开锁日志发生异常", err)
		return
	}
	logs.Debug("按键开锁日志记录成功", "DeviceId:", deviceId)
}

// 记录用户开锁日志
func UserUnlockLog(deviceId string, userName string, method string, success bool) {
	var log models.UserUnlockLog
	log.UserName = userName
	log.DeviceId = deviceId
	log.Method = method
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
	logs.Debug("用户开锁日志记录成功", "DeviceId:", deviceId, "UserName:", userName, "Method:", method)

}
