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
func SendGet(remote models.ConfigRemote) (request interface{}, err error) {

	url := remote.RemoteAddr
	requ := httplib.Get(url)

	err = setReqParam(requ, remote)
	if err != nil {
		return request, err
	}
	return requ, nil
}

// post 请求
func SendPost(remote models.ConfigRemote) (request interface{}, err error) {

	url := remote.RemoteAddr
	requ := httplib.Post(url)
	err = setReqParam(requ, remote)

	if err != nil {
		return request, err
	}
	return requ, nil
}

// put 请求
func SendPut(remote models.ConfigRemote) (request interface{}, err error) {

	url := remote.RemoteAddr
	requ := httplib.Put(url)
	err = setReqParam(requ, remote)

	if err != nil {
		return request, err
	}
	return requ, nil
}

// delete请求
func SendDelete(remote models.ConfigRemote) (request interface{}, err error) {

	url := remote.RemoteAddr
	requ := httplib.Delete(url)
	err = setReqParam(requ, remote)

	if err != nil {
		return request, err
	}
	return requ, nil
}

// header请求
func SendHeader(remote models.ConfigRemote) (request interface{}, err error) {

	url := remote.RemoteAddr
	requ := httplib.Head(url)
	err = setReqParam(requ, remote)

	if err != nil {
		return request, err
	}
	return requ, nil
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
