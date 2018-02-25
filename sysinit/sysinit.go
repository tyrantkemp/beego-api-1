package sysinit

import (
	"github.com/astaxie/beego"
	"beego-api-1/utils"
)

func init(){


	// 开启session
	beego.BConfig.WebConfig.Session.SessionOn=true
	// 初始化日志功能
	utils.InitLogs()
	// 初始化缓存功能
	utils.InitCache()
	// 初始化数据库
	InitDataBase()
}
