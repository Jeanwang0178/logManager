package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"],
		beego.ControllerComments{
			Method:           "Edit",
			Router:           `/edit`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/list`,
			AllowHTTPMethods: []string{"post", "get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"],
		beego.ControllerComments{
			Method:           "Save",
			Router:           `/save`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:BizLogController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/view`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ConfigController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ConfigController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/view`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ConfigController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ConfigController"],
		beego.ControllerComments{
			Method:           "Write",
			Router:           `/write`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ContentController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ContentController"],
		beego.ControllerComments{
			Method:           "ListRemoteFile",
			Router:           `/listRemoteFile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ContentController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ContentController"],
		beego.ControllerComments{
			Method:           "QueryContent",
			Router:           `/queryContent`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ContentController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ContentController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/view`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/delete`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"],
		beego.ControllerComments{
			Method:           "Edit",
			Router:           `/edit`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get", "post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"],
		beego.ControllerComments{
			Method:           "Save",
			Router:           `/save`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DatabaseController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/view`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"],
		beego.ControllerComments{
			Method:           "GetTime",
			Router:           `/gettime`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"get", "post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"],
		beego.ControllerComments{
			Method:           "Logout",
			Router:           `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:DefaultController"],
		beego.ControllerComments{
			Method:           "Profile",
			Router:           `/profile`,
			AllowHTTPMethods: []string{"get", "post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:FieldController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:FieldController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:FieldController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:FieldController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:FieldController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:FieldController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:FieldController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:FieldController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:FieldController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:FieldController"],
		beego.ControllerComments{
			Method:           "Edit",
			Router:           `/edit`,
			AllowHTTPMethods: []string{"get", "post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:FieldController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:FieldController"],
		beego.ControllerComments{
			Method:           "Save",
			Router:           `/save`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:HelpController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:HelpController"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"],
		beego.ControllerComments{
			Method:           "ListRemoteFile",
			Router:           `/listRemoteFile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"],
		beego.ControllerComments{
			Method:           "TailfLog",
			Router:           `/tailfLog`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/view`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:LogFileController"],
		beego.ControllerComments{
			Method:           "ViewLog",
			Router:           `/viewLog`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ManagerController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ManagerController"],
		beego.ControllerComments{
			Method:           "DataList",
			Router:           `/dataList`,
			AllowHTTPMethods: []string{"get", "post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ManagerController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ManagerController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:ManagerController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:ManagerController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/view`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:MonitorController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:MonitorController"],
		beego.ControllerComments{
			Method:           "ListRemoteFile",
			Router:           `/listRemoteFile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:MonitorController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:MonitorController"],
		beego.ControllerComments{
			Method:           "QueryContent",
			Router:           `/queryContent`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:MonitorController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:MonitorController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/view`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/delete`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"],
		beego.ControllerComments{
			Method:           "Edit",
			Router:           `/edit`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"],
		beego.ControllerComments{
			Method:           "KafkaList",
			Router:           `/kafkaList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"],
		beego.ControllerComments{
			Method:           "Save",
			Router:           `/save`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"] = append(beego.GlobalControllerRouter["logManager/src/controllers:RemoteController"],
		beego.ControllerComments{
			Method:           "SaveAddr",
			Router:           `/saveAddr`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

}
