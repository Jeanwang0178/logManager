package controllers

//HelpController
type HelpController struct {
	BaseController
}

// @router / [get]
func (ctl *HelpController) Index() {
	ctl.Data["pageTitle"] = "使用帮助"
	ctl.display()
}
