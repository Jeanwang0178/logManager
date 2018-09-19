package controllers

import (
	"github.com/astaxie/beego"
	"time"
)

type DefaultController struct {
	BaseController
}

// @router / [get]
func (ctl *DefaultController) Index() {
	ctl.Data["Website"] = "beego.me"
	ctl.Data["Email"] = "wjian0124@163.com"
	ctl.TplName = "index.html"
	ctl.display()
}

// 获取系统时间
// @router /gettime [get]
func (this *DefaultController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}

// 退出登录
func (this *DefaultController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor("MainController.Login"))
}
