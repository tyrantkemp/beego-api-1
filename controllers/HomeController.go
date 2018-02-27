package controllers

import (
	"beego-api-1/models"
	"encoding/json"
	"beego-api-1/enums"
	"beego-api-1/utils"
	"github.com/astaxie/beego"
)

type HomeController struct {
	BaseController
}

// URLMapping ...
func (c *HomeController) URLMapping() {

}

//@Title login
//@Description  login
//@Param body body models.Login true  "body for login"
//@Success  {object} Common.JsonResult
//@Failure  {object} common.JsonResult
//@router /login [post]
func (c *HomeController) Login() {
	var login models.Login
	json.Unmarshal(c.Ctx.Input.RequestBody, &login)
	username := login.UserName
	userpassword := login.UserPwd
	if len(username) == 0 || len(userpassword) == 0 {
		c.jsonResult(enums.JRCodeFail, "用户名和密码不正确", "")
	}
	userpwd := utils.String2md5(userpassword)
	user, err := models.GetUserByNameAndPwd(username, userpwd)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.jsonResult(enums.JRCodeFail, "用户状态不可用", "")
		}
		//TODO 用户信息及所权限信息保存在redis里
		// 生成用户唯一token
		token := "TOKEN_" + utils.UniqueId()
		// 保存个人信息到redis中，设置过期时间为48小时
		expireTime, err := beego.AppConfig.Int("tokenExpireTime")
		if err != nil {
			c.jsonResult(enums.JRCodeFail, "token过期时间设置错误", "")
		}
		c.setUserToRedis(token, user.Id, expireTime)
		//返回用户token
		c.jsonResult(enums.JRCodeSuccess, "用户登录成功", token)
	} else {
		c.jsonResult(enums.JRCodeFail, "用户名或密码错误", "")

	}

}

//@router /index [post]
func (c *HomeController) Index() {
	utils.LogInfo("token:" + c.Ctx.Request.Header.Get("token"))

}

//@Title logout
//@Description  logout
//@Param token header string true  "token of user"
//@Success  {object} Common.JsonResult
//@Failure  {object} common.JsonResult
// @router /logout [post]
func (c *HomeController) Logout() {
	utils.LogInfo("退出登录...")
	token := c.Ctx.Request.Header.Get("token")
	// value 为空 底层操作为直接删除该token对应数据
	user := models.User{}
	err := utils.SetCache(token, user, 1)
	if err != nil {
		c.jsonResult(enums.JRCodeFail, "用户登出失败", err)
	}
	c.jsonResult(enums.JRCodeSuccess, "用户登出成功", "")
}
