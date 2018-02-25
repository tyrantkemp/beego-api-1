package utils

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"strings"
)

var consolelogs *logs.BeeLogger
var filelogs *logs.BeeLogger

var runmode string

func InitLogs() {

	// 开发环境
	consolelogs = logs.NewLogger(1)
	consolelogs.SetLogger(logs.AdapterConsole)
	consolelogs.Async()
	consolelogs.EnableFuncCallDepth(true)
	consolelogs.SetLogFuncCallDepth(4)

	// 生产环境
	filelogs = logs.NewLogger(10000)
	level := beego.AppConfig.String("log::level")
	//
	filelogs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/api.log",
		"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
		"level":`+level+`,
		"daily":true,
		"maxdays":10}`)

	filelogs.Async()
	filelogs.EnableFuncCallDepth(true)
	filelogs.SetLogFuncCallDepth(4)
	runmode = strings.TrimSpace(strings.ToLower(beego.AppConfig.String("runmode")))
	if runmode == "" {
		runmode = "dev"
	}

}

func LogEmergency(v interface{}) {
	log("emergency", v)
}

func LogAlert(v interface{}) {

	log("alert", v)

}

func LogCritical(v interface{}) {
	log("critical", v)
}

func LogError(v interface{}) {

	log("error", v)

}

func LogWarning(v interface{}) {

	log("warning", v)

}

func LogNotice(v interface{}) {
	log("notice", v)
}

func LogInfo(v interface{}) {
	log("info", v)
}
func LogDebug(v interface{}) {
	log("debug", v)
}
func LogTrace(v interface{}) {
	log("trace", v)
}

func log(level, v interface{}) {
	formate := "%s"

	if level == "" {
		level = "debug"
	}

	if runmode == "dev" {
		switch level {
		case "emergency":
			consolelogs.Emergency(formate, v)
		case "alert":
			consolelogs.Alert(formate, v)
		case "critical":
			consolelogs.Critical(formate, v)
		case "error":
			consolelogs.Error(formate, v)
		case "warning":
			consolelogs.Warning(formate, v)
		case "notice":
			consolelogs.Notice(formate, v)
		case "info":
			consolelogs.Info(formate, v)
		case "debug":
			consolelogs.Debug(formate, v)
		case "trace":
			consolelogs.Trace(formate, v)
		default:
			consolelogs.Debug(formate, v)
		}
	} else if runmode == "prod" {

		switch level {
		case "emergency":
			filelogs.Emergency(formate, v)
		case "alert":
			filelogs.Alert(formate, v)
		case "ctritical":
			filelogs.Critical(formate, v)
		case "error":
			filelogs.Error(formate, v)
		case "warning":
			filelogs.Warning(formate, v)
		case "notice":
			filelogs.Notice(formate, v)
		case "info":
			filelogs.Info(formate, v)
		case "debug":
			filelogs.Debug(formate, v)
		case "trace":
			filelogs.Trace(formate, v)
		default:
			filelogs.Debug(formate, v)
		}
	}

}
