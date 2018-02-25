package controllers

import (
	"github.com/astaxie/beego"
	"beego-api-1/utils"
	"fmt"
)

type HelloController struct {
	beego.Controller
}

//@router /put [get]
func (H* HelloController)Put(){
	err:=utils.SetCache("test-name","xiaozhun",30)
	if err!=nil{
		utils.LogError(err)
	}

}

//@router /get [get]
func (H* HelloController)Get(){
	var name string
 	err:=utils.GetCache("test-name",&name)
	if err!=nil{
		utils.LogError(err)
	}
	utils.LogInfo("get cache value:"+name)
	fmt.Println("name",name)
}
