package controllers

import (
	"bufio"
	"bytes"
	"github.com/axgle/mahonia"
	"logManager/src/common"
	"logManager/src/utils"
	"os"
	"path/filepath"
)

type ConfigController struct {
	BaseController
}

// @router /view [get]
func (ctl *ConfigController) View() {
	response := make(map[string]interface{})
	filePath1, _ := filepath.Abs("./") //C:\goWorkSpace\src\logManager
	common.Logger.Info(filePath1)
	confStr, err := utils.ReadFile("conf/app.conf")
	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = confStr
	}
	ctl.Data["result"] = response
	ctl.display()
}

// @router /write [post]
func (ctl *ConfigController) Write() {
	response := make(map[string]interface{})
	content := ctl.GetString("content")
	data := []byte(content)
	err := utils.WriteFile(data, "C:/goWorkSpace/src/logManager/conf/app.conf")
	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
	}
	ctl.Data["json"] = response
	ctl.ServeJSON()
}

func bufioRead(fileName string) (confStr string, err error) {
	fileObj, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer fileObj.Close()
	decoder := mahonia.NewDecoder("utf-8")
	reader := bufio.NewReader(fileObj)
	buf := make([]byte, 1024)
	var result bytes.Buffer

	//读取Reader对象中的内容到[]byte类型的buf中
	if _, err := reader.Read(buf); err == nil {
		result.WriteString(decoder.ConvertString(string(buf)))
	}
	return result.String(), nil

}
