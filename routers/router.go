package routers

import (
	"github.com/astaxie/beego"
	"logManager/src/controllers"
	_ "logManager/src/models"
)

func init() {

	ns := beego.NewNamespace("/open",
		beego.NSNamespace("/index",
			beego.NSInclude(
				&controllers.DefaultController{},
			),
		),
		beego.NSNamespace("/bizLog",
			beego.NSInclude(
				&controllers.BizLogController{},
			),
		),
		beego.NSNamespace("/manager",
			beego.NSInclude(
				&controllers.ManagerController{},
			),
		),
		beego.NSNamespace("/content",
			beego.NSInclude(
				&controllers.ContentController{},
			),
		),
		beego.NSNamespace("/logFile",
			beego.NSInclude(
				&controllers.LogFileController{},
			),
		),
		beego.NSNamespace("/config",
			beego.NSInclude(
				&controllers.ConfigController{},
			),
		),
		beego.NSNamespace("/dataBase",
			beego.NSInclude(
				&controllers.DatabaseController{},
			),
		),
		beego.NSNamespace("/field",
			beego.NSInclude(
				&controllers.FieldController{},
			),
		),
		beego.NSNamespace("/remote",
			beego.NSInclude(
				&controllers.RemoteController{},
			),
		),

		beego.NSNamespace("/help",
			beego.NSInclude(
				&controllers.HelpController{},
			),
		),
	)

	beego.AddNamespace(ns)

}
