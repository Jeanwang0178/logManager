package controllers

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
