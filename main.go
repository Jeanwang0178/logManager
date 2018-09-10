package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/bee/generate/swaggergen"
	"github.com/dwdcth/consoleEx"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"html/template"
	"logManager/inital"
	_ "logManager/routers"
	"logManager/utils"
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

	out := consoleEx.ConsoleWriterEx{Out: colorable.NewColorableStdout()}

	zerolog.CallerSkipFrameCount = 2 // 根据实际，另外获取的是MSG调用处的文件路径和行号

	logger := zerolog.New(out).With().Caller().Timestamp().Logger()

	logger.Info().Msg("info")

	logger.Debug().Msg("debug")

	beego.AddFuncMap("GetMapValue", utils.GetMapValue)
	beego.AddFuncMap("GetSliceMapValue", utils.GetSliceMapValue)

	swaggergen.GenerateDocs("")

	beego.Run()
}
