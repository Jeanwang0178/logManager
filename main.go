package main

import (
	"github.com/astaxie/beego"
	"html/template"
	"logManager/inital"
	_ "logManager/routers"
	_ "logManager/utils"
	"net/http"
)

const VERSION = "1.0.0"

func main() {

	inital.Init()

	//设置默认404页面
	beego.ErrorHandler("404", func(writer http.ResponseWriter, request *http.Request) {
		t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
		data := make(map[string]interface{})
		data["content"] = "page not found"
		t.Execute(writer, data)
	})

	beego.AppConfig.Set("version", VERSION)
	beego.BConfig.WebConfig.Session.SessionOn = true

	beego.Run()
}
