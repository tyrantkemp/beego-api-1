package controllers

import (
	"github.com/astaxie/beego"
	"beego-api-1/models"
	"beego-api-1/enums"
	"strings"
	"beego-api-1/utils"
	"fmt"
)

type BaseController struct {
	beego.Controller

	controllerName string
	actionName     string
	curUser        models.User
}

func (c *BaseController) Prepare() {
	c.controllerName, c.actionName = c.GetControllerAndAction()
	c.adapterUserInfo()
}

//从session里取用户信息
func (c *BaseController) adapterUserInfo() {

	// session
	/*a := c.GetSession(enums.CURRENT_USER)
	if a != nil {
		c.curUser = a.(models.User)
		c.Data[enums.CURRENT_USER] = a
	}*/

	//redis
	token := c.Ctx.Request.Header.Get("token")
	var user models.User
	utils.GetCache(token, &user)
/*	if err != nil {
		c.jsonResult(enums.JRCodeFail, "从redis获取用户信息失败", err)
	}*/
	if &user != nil {
		c.curUser = user
	}

}

// 将用户信息保存到session 包括拥有的权限
func (c *BaseController) setUserToSession(userId int) error {
	user, err := models.GetUserById(userId)
	if err != nil {
		return err
	}
	resouceUrls := models.GetResouceUrlByUserId(userId)
	for _, item := range resouceUrls {
		user.ResourceUrlList = append(user.ResourceUrlList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession(enums.CURRENT_USER, *user)
	return nil
}

//
func (c *BaseController) setUserToRedis(token string, userId int, expireTime int) error {
	user, err := models.GetUserById(userId)
	if err != nil {
		return err
	}
	resouceUrls := models.GetResouceUrlByUserId(userId)
	for _, item := range resouceUrls {
		user.ResourceUrlList = append(user.ResourceUrlList, strings.TrimSpace(item.UrlFor))
	}
	utils.SetCache(token, user, expireTime)
	return nil

}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		//登录页面地址
		//urlstr := c.URLFor("HomeController.Login") + "?url="
		//登录成功后返回的址为当前
	//	returnURL := c.Ctx.Request.URL.Path
		//如果ajax请求则返回相应的错码和跳转的地址
		/*if c.Ctx.Input.IsAjax() {
			//由于是ajax请求，因此地址是header里的Referer
			returnURL := c.Ctx.Input.Refer()
			c.jsonResult(enums.JRCode302, "请登录", urlstr+returnURL)
		}
		c.Redirect(urlstr+returnURL, 302)*/
		c.jsonResult(enums.JRCodeFail,"请登录","")
		c.StopRun()
	}

}

// 判断当前用户是否有访问 controller.action 的权限
func (c *BaseController) checkActionPermission(controller string, action string) bool {

	if c.curUser.Id == 0 {
		return false
	}
	//session
	//user := c.GetSession(enums.CURRENT_USER)

	//redis

	// 类型断言
	//v, ok := user.(models.User)
	//if ok {
	v := c.curUser
	// 如果超级管理员
	if v.IsAdmin == true {
		return true
	}

	for i, _ := range v.ResourceUrlList {
		urlfor := strings.TrimSpace(v.ResourceUrlList[i])
		if len(urlfor) == 0 {
			continue
		}
		strs := strings.Split(urlfor, ",")
		if len(strs[0]) > 0 && strs[0] == (controller+"."+action) {
			return true
		}
	}
	//	}
	return false

}

func (c *BaseController) checkPermission(ignores ... string) {

	c.checkLogin()

	for _, action := range ignores {
		if action == c.actionName {
			return
		}
	}

	isPermitted := c.checkActionPermission(c.controllerName, c.actionName)
	if !isPermitted {

		utils.LogDebug(fmt.Sprintf("author controll path:%s.%s userId:%v 无权限访问 ", c.controllerName, c.actionName, c.curUser.Id))
		if c.Ctx.Input.IsAjax() {
			c.jsonResult(enums.JRCode401, "无权限访问", "")

		} else {

			c.pageError("无权限访问")
		}
	}

}
func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{Code: code, Msg: msg, Obj: obj}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}

// 重定向去错误页h面
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("HomeController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()

}

// 重定向到登录页
func (c *BaseController) pageLogin() {
	url := c.URLFor("HomeController.Login")
	c.Redirect(url, 302)
	c.StopRun()
}

// 重定向
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}
