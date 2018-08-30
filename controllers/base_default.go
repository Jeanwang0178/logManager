package controllers

type DefaultController struct {
	BaseController
}

// @router / [get]
func (ctl *DefaultController) Index() {
	ctl.Data["Website"] = "beego.me"
	ctl.Data["Email"] = "astaxie@gmail.com"
	ctl.TplName = "index.html"
	ctl.display()
}
