package controllers

import (
	"bufio"
	"bytes"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"logManager/common"
	"logManager/utils"
	"os"
)

type ConfigController struct {
	BaseController
}

// @router /view [get]
func (ctl *ConfigController) View() {
	response := make(map[string]interface{})
	//confStr, err := bufioRead("C:/goWorkSpace/src/logManager/conf/app.conf")
	confStr, err := ReadFile("C:/goWorkSpace/src/logManager/conf/app.conf")
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
	err := writeFile(data, "C:/goWorkSpace/src/logManager/conf/app.conf")
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
	common.Logger.Info(result.String())
	return result.String(), nil

}

func ReadFile(fileName string) (st string, err error) {
	//ret := make([]string, 0)
	var result bytes.Buffer
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		result.WriteString(scanner.Text() + "\n")
		//ret = append(ret, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	common.Logger.Info(result.String())
	return result.String(), nil
}

func writeFile(content []byte, fileName string) error {

	err := ioutil.WriteFile(fileName, content, 0666)
	if err != nil {
		return err
	}
	return nil
}
