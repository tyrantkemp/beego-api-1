package sysinit

import (
	_ "beego-api-1/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	_ "github.com/go-sql-driver/mysql"
)

func InitDataBase() {

	dbName := beego.AppConfig.String("db::name")
	dbUser := beego.AppConfig.String("db::user")
	dbPwd := beego.AppConfig.String("db::pwd")
	dbHost := beego.AppConfig.String("db::host")
	dbPort := beego.AppConfig.String("db::port")

	orm.RegisterDataBase("default", "mysql", dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":"+
		dbPort+ ")/"+ dbName+ "?charset=utf8", 30)

	isDev := (beego.AppConfig.String("runmode") == "dev")

	orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}

}
