package routers

import (
	"github.com/astaxie/beego"
	"logManager/controllers"
)

func init() {

	ns := beego.NewNamespace("/open",
		beego.NSNamespace("/index",
			beego.NSInclude(
				&controllers.DefaultController{},
			),
		),
		beego.NSNamespace("/logManager",
			beego.NSInclude(
				&controllers.BizLogController{},
			),
		),
		beego.NSNamespace("/log",
			beego.NSInclude(
				&controllers.ManagerController{},
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
		beego.NSNamespace("/help",
			beego.NSInclude(
				&controllers.HelpController{},
			),
		),
	)

	beego.AddNamespace(ns)

}
