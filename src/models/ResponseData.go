package models

type ResponseData struct {
	Code string        `json:"code"`
	Data []interface{} `json:"data"`
	Msg  string        `json:"msg"`
}

type RequestFileParam struct {
	RemoteAddr  string `json:"remoteAddr"`
	FilePath    string `json:"filePath"`
	Content     string `json:"content"`
	Position    int    `json:"position"`
	PreLineNum  int64  `json:"preLineNum"`
	NextLineNum int64  `json:"nextLineNum"`
	LineNum     int64  `json:"lineNum"`
	QueryType   string `json:"queryType"`
	OperType    string `json:"operType"`
	QueryCount  int    `json:"queryCount"`
}
