package models

type ResponseData struct {
	Code string        `json:"code"`
	Data []interface{} `json:"data"`
	Msg  string        `json:"msg"`
}
