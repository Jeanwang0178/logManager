package services

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/eapache/queue"
	"io"
	"logManager/src/common"
	"logManager/src/models"
	"os"
	"strings"
)

const (
	queryOffSet = int64(5120)
)

/**
 * 搜索、查询文件内容
 */
func MonitorFileServiceQueryContent(param models.RequestFileParam) (data interface{}, preOffset int64, nextOffset int64, retPosition int, err error) {

	remoteAddr := strings.TrimSpace(param.RemoteAddr)
	common.Logger.Info(remoteAddr)
	filePath := strings.TrimSpace(param.FilePath)
	content := strings.TrimSpace(param.Content)
	PreLineNum := param.PreLineNum
	NextLineNum := param.NextLineNum
	QueryType := param.QueryType
	OperType := param.OperType
	Position := param.Position

	isSearch := false
	//1、查询方式: 搜索 FALSE 或 定位
	if !(OperType == "scroll" || strings.TrimSpace(content) == "") {
		isSearch = true
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", PreLineNum, NextLineNum, Position, err
	}
	defer file.Close()

	toal, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, PreLineNum, NextLineNum, Position, err
	}

	//2、定位文件
	location := int64(0)
	location, err = getFileSeekLocation(isSearch, QueryType, location, PreLineNum, NextLineNum, toal)
	if err != nil {
		return nil, PreLineNum, NextLineNum, Position, err
	}
	preOffset = location
	offset, err := file.Seek(location, io.SeekStart)
	if err != nil {
		return nil, PreLineNum, NextLineNum, Position, err
	}
	common.Logger.Info("start location ============%d =========", offset)

	//3、搜索函数 /定位函数
	result := make(map[string]interface{})
	if isSearch {
		searchFileContent(content, file, Position, location, toal, QueryType, result)
	} else {
		queryFileContent(file, location, toal, result)
	}

	if result["err"] != nil {
		err = result["err"].(error)
		if err != nil {
			if !strings.Contains(err.Error(), "1001") {
				common.Logger.Error("search is error :", err)
				return nil, PreLineNum, NextLineNum, Position, err
			}
		}
	}
	preQueue := result["retQueue"].(*queue.Queue)
	preOff := result["preOffSet"].(int64)
	nextOff := result["nextOffset"].(int64)
	if result["position"] == nil {

	} else {
		retPosition = result["position"].(int)
	}

	if !isSearch && QueryType == "N" {
		preOffset = PreLineNum
	} else {
		preOffset = preOff
	}

	if !isSearch && QueryType == "P" {
		nextOffset = NextLineNum
	} else {
		nextOffset = nextOff
	}

	common.Logger.Info("===============pre = %d ======next == %d =========", preOffset, nextOffset)

	//5、返回查询字符串 、定位：头、尾 、err
	var buffer = bytes.Buffer{}
	for i := 0; i < preQueue.Length(); i++ {
		buffer.WriteString(preQueue.Get(i).(string))
	}

	return buffer.String(), preOffset, nextOffset, retPosition, err
}

//搜索文件内容
func searchFileContent(search string, file *os.File, position int, offset int64, toal int64, queryType string, result map[string]interface{}) {
	oldPosition := position
	if "P" == queryType {
		if position >= 1 {
			position -= 1
		} else {
			position = 0
		}
	} else {
		position += 1
	}

	cmp := []byte(search)
	retQueue := queue.New()
	findCount := 0
	curOffset := int64(0)
	input := bufio.NewScanner(file)

	isEnd := false
	preOffSet := int64(0)
	nextOffset := int64(0)
	findOffset := int64(0)

	for input.Scan() { //搜索内容
		info := input.Bytes()
		curOffset += int64(len(info)) + 1 //1 代表换行
		if findCount < position {
			if bytes.Contains(info, cmp) {
				findCount++
			}
			if findCount == position {
				findOffset = curOffset
				curOffset = 0
			} else {
				continue
			}
		}
		//搜索到内容
		stext := string(info)
		retQueue.Add(stext + "\n")

		if curOffset >= queryOffSet {
			preOffSet = offset + findOffset
			nextOffset = curOffset + findOffset + offset
			isEnd = true
			break
		}
	}

	if findCount >= position { // 查询到内容
		common.Logger.Info("========offset========%d ===findCount=====%d ====position=====%d=======", offset, findCount, position)

		if !isEnd { //扫描所有文件内容
			preOffSet = offset + findOffset
			nextOffset = toal
		}
		result["preOffSet"] = preOffSet
		result["nextOffset"] = nextOffset
		result["retQueue"] = retQueue
		result["position"] = position

	} else { //未查询到内容

		repeat := int64(0)
		if "P" == queryType { //已经查询到文件头部
			repeat = 0
			common.Logger.Info("已经查询到文件头部,截取顶部内容")

		} else { //已经查询到文件尾部
			if toal-queryOffSet > 0 {
				repeat = toal - queryOffSet
			} else {
				repeat = 0
			}
			common.Logger.Info("已经查询到文件尾部,截取结尾内容")
		}
		file.Seek(repeat, io.SeekStart) //获取文件头部或/尾部内容
		err := errors.New("1001:" + "已搜索全部内容,没有匹配的关键词【" + search + "】")
		result["err"] = err
		result["position"] = oldPosition
		queryFileContent(file, repeat, toal, result)
	}

}

// 查询文件内容（非搜索）

func queryFileContent(file *os.File, offset int64, toal int64, result map[string]interface{}) {

	retQueue := queue.New()
	curOffset := int64(0)
	input := bufio.NewScanner(file)

	isEnd := false
	preOffSet := int64(0)
	nextOffset := int64(0)

	for input.Scan() {
		info := input.Bytes()
		curOffset += int64(len(info)) + 1 //1 代表换行

		stext := string(info)
		retQueue.Add(stext + "\n")

		if curOffset >= queryOffSet {
			preOffSet = offset
			nextOffset = curOffset + offset
			isEnd = true
			break
		}
	}

	if !isEnd { //扫描所有文件内容
		preOffSet = offset
		nextOffset = toal
	}
	result["preOffSet"] = preOffSet
	result["nextOffset"] = nextOffset
	result["retQueue"] = retQueue

}

/** 定位文件位置
 */
func getFileSeekLocation(isSearch bool, QueryType string, location int64, PreLineNum int64, NextLineNum int64, toal int64) (int64, error) {
	if isSearch {
		if "N" == QueryType {
			location = NextLineNum //向下
		} else if "P" == QueryType {
			if PreLineNum > queryOffSet {
				location = PreLineNum - queryOffSet //分段向上
			} else {
				location = 0
			}
		}
	} else {
		if "N" == QueryType {
			if NextLineNum == toal {
				common.Logger.Info("The end of the file has been searched")
				return 0, errors.New("1001:" + "已经搜索到文件结尾")
			}
			location = NextLineNum //向下
		} else if "P" == QueryType {
			if PreLineNum == 0 {
				common.Logger.Info("The end of the file has been searched")
				return 0, errors.New("1001:" + "已经搜索到文件开头")
			}
			if PreLineNum > queryOffSet {
				location = PreLineNum - queryOffSet //分段向上
			} else {
				location = 0
			}
		} else if "T" == QueryType {
			if toal-queryOffSet > 0 {
				location = toal - queryOffSet //底部向上
			} else {
				location = 0
			}
		} else if "H" == QueryType {
			location = 0
		}
	}
	return location, nil
}
