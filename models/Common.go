package models

import "beego-api-1/enums"

type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg string `json:"msg"`
	Obj interface{} `json:"obj"`
}


