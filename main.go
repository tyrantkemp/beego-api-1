package main

import (
	_ "beego-api-1/routers"

	"github.com/astaxie/beego"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"beego-api-1/controllers"

	_ "beego-api-1/sysinit"
)

func main() {
	// 解决swagger跨域访问
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// josnRpc
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(controllers.HelloService), "")
	beego.Handler("/rpc", s)

	beego.Run()
}
