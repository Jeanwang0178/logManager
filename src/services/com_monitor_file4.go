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
	queryOffSet4 = int64(5120)
)

/**
 * 1\出现次数
 */
func MonitorFileServiceQueryContent4(param models.RequestFileParam) (data interface{}, preOffset int64, nextOffset int64, position int, err error) {

	remoteAddr := strings.TrimSpace(param.RemoteAddr)
	common.Logger.Info(remoteAddr)
	filePath := strings.TrimSpace(param.FilePath)
	content := strings.TrimSpace(param.Content)
	PreLineNum := param.PreLineNum
	NextLineNum := param.NextLineNum
	QueryType := param.QueryType
	OperType := param.OperType
	QueryCount := param.QueryCount

	isSearch := false
	//1、查询方式: 搜索 FALSE 或 定位
	if !(OperType == "scroll" || strings.TrimSpace(content) == "") {
		isSearch = true
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", PreLineNum, NextLineNum, position, err
	}

	toal, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, PreLineNum, NextLineNum, position, err
	}

	//2、定位文件
	location := int64(0)
	location, err = getFileSeekLocation4(isSearch, QueryType, location, PreLineNum, NextLineNum, toal)
	if err != nil {
		return nil, PreLineNum, NextLineNum, position, err
	}
	preOffset = location
	offset, err := file.Seek(location, io.SeekStart)
	if err != nil {
		return nil, PreLineNum, NextLineNum, position, err
	}
	common.Logger.Info("start location ============%d =========", offset)
	defer file.Close()

	//3、搜索函数 /定位函数
	result := make(map[string]interface{})
	searchFileContent4(isSearch, content, file, QueryCount, location, toal, QueryType, result, true)
	if result["err"] != nil {
		err = result["err"].(error)
		if err != nil {
			if !strings.Contains(err.Error(), "1001") {
				common.Logger.Error("search is error :", err)
				return nil, PreLineNum, NextLineNum, position, err
			}
		}
	}
	preQueue := result["retQueue"].(*queue.Queue)
	preOff := result["preOffSet"].(int64)
	nextOff := result["nextOffset"].(int64)
	position = result["position"].(int)

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

	return buffer.String(), preOffset, nextOffset, position, err
}

func getFileSeekLocation4(isSearch bool, QueryType string, location int64, PreLineNum int64, NextLineNum int64, toal int64) (int64, error) {
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

func searchFileContent4(isSearch bool, search string, file *os.File, position int, offset int64, toal int64, queryType string, result map[string]interface{}, isQuery bool) {

	cmp := []byte(search)
	retQueue := queue.New()
	findCount := 0
	curOffset := int64(0)
	input := bufio.NewScanner(file)

	isEnd := false
	preOffSet := int64(0)
	nextOffset := int64(0)
	findOffset := int64(0) //记录查询首次位置

	for input.Scan() {
		info := input.Bytes()
		curOffset += int64(len(info)) + 1 //1 代表换行
		if isSearch && isQuery && findCount <= position {
			if bytes.Contains(info, cmp) {
				findCount++
				findOffset += curOffset
				curOffset = 0
			}
			continue
		}
		stext := string(info)
		retQueue.Add(stext + "\n")

		if curOffset >= queryOffSet {
			if findCount > position {
				preOffSet = offset + findOffset
				nextOffset = curOffset + findOffset + offset
			} else {
				preOffSet = offset
				nextOffset = curOffset + offset
			}
			isEnd = true
			break
		}
	}

	if isSearch && findCount <= position && isQuery { //isAll 重新查找内容
		repeat := int64(0)
		if "P" == queryType {
			if offset > 0 { //未查询到文件头部，继续查找
				preOffSet = 0
				if offset-queryOffSet4 > 0 {
					repeat = offset - queryOffSet4
				} else {
					repeat = 0
				}
				isQuery = true
			} else {
				repeat = 0
				isQuery = false
			}

		} else { //已经查询到文件尾部
			if toal-queryOffSet4 > 0 {
				repeat = toal - queryOffSet4
			} else {
				repeat = 0
			}
			isQuery = false
		}
		file.Seek(repeat, io.SeekStart) //获取文件头部或/尾部内容
		err := errors.New("1001:" + "没有匹配的关键词【" + search + "】")
		result["err"] = err
		searchFileContent4(isSearch, search, file, findCount, repeat, toal, queryType, result, isQuery)
	} else {
		common.Logger.Info("===================offset===========%d ===findCount=======%d =====", offset, findCount)

		if !isEnd { //扫描所有文件内容
			preOffSet = offset
			nextOffset = toal
		}
		result["preOffSet"] = preOffSet
		result["nextOffset"] = nextOffset
		result["retQueue"] = retQueue
		result["position"] = position
	}

}
