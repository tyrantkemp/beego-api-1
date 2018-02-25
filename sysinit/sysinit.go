package sysinit

import (
	"github.com/astaxie/beego"
	"beego-api-1/utils"
)

func init(){


	beego.BConfig.WebConfig.Session.SessionOn=true
	utils.InitLogs()
	utils.InitCache()
	InitDataBase()
}
