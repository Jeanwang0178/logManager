package utils

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"logManager/models"
)

const (
	pheader = "header"
	pparam  = "param"
	pbody   = "body"
)

// get 请求
func SendGet(remote models.ConfigRemote) (resp string, err error) {

	url := remote.RemoteAddr
	req := httplib.Get(url)
	err = setReqParam(req, remote)
	if err != nil {
		return "", err
	}

	resp, err = req.String()
	if err != nil {
		return "", err
	}
	return resp, nil
}

// post 请求
func SendPost(remote models.ConfigRemote) (resp string, err error) {

	url := remote.RemoteAddr
	req := httplib.Post(url)
	err = setReqParam(req, remote)

	resp, err = req.String()

	if err != nil {
		return "", err
	}
	return resp, nil
}

// put 请求
func SendPut(remote models.ConfigRemote) (resp string, err error) {

	url := remote.RemoteAddr
	req := httplib.Put(url)
	err = setReqParam(req, remote)

	resp, err = req.String()
	if err != nil {
		return "", err
	}
	return resp, nil
}

// delete请求
func SendDelete(remote models.ConfigRemote) (resp string, err error) {

	url := remote.RemoteAddr
	req := httplib.Delete(url)
	err = setReqParam(req, remote)

	resp, err = req.String()
	if err != nil {
		return "", err
	}
	return resp, nil
}

// header请求
func SendHeader(remote models.ConfigRemote) (resp string, err error) {

	url := remote.RemoteAddr
	req := httplib.Head(url)
	err = setReqParam(req, remote)

	resp, err = req.String()
	if err != nil {
		return "", err
	}
	return resp, nil
}

//设置请求参数
func setReqParam(req *httplib.BeegoHTTPRequest, remote models.ConfigRemote) (err error) {

	jsonHeader := remote.Header
	jsonParam := remote.Param
	jsonBody := remote.Body

	if jsonHeader == "" {
		jsonHeader = "{}"
	}

	if jsonParam == "" {
		jsonParam = "{}"
	}

	if jsonBody == "" {
		jsonBody = "{}"
	}

	// request header
	jsonByte := []byte(jsonHeader)
	headerMap := make(map[string]interface{})
	err = json.Unmarshal(jsonByte, &headerMap)
	if err != nil {
		return err
	} else {
		for key, value := range headerMap {
			req.Header(key, value.(string))
		}
	}

	// request param
	jsonByte = []byte(jsonParam)
	paramMap := make(map[string]interface{})
	err = json.Unmarshal(jsonByte, &paramMap)
	if err != nil {
		return err
	} else {
		for key, value := range paramMap {
			req.Param(key, value.(string))
		}
	}

	// request body
	jsonByte = []byte(jsonBody)
	bodyMap := make(map[string]interface{})
	err = json.Unmarshal(jsonByte, &bodyMap)
	if err != nil {
		return err
	} else {
		req.JSONBody(bodyMap)
	}

	return nil
}
