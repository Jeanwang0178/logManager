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
	queryOffSet2 = int64(5120)
)

/**
 * 1\出现次数
 */
func MonitorFileServiceQueryContent2(param models.RequestFileParam) (data interface{}, preOffset int64, nextOffset int64, err error) {

	remoteAddr := strings.TrimSpace(param.RemoteAddr)
	common.Logger.Info(remoteAddr)
	filePath := strings.TrimSpace(param.FilePath)
	content := strings.TrimSpace(param.Content)
	PreLineNum := param.PreLineNum
	NextLineNum := param.NextLineNum
	QueryType := param.QueryType
	OperType := param.OperType

	isSearch := false
	//1、查询方式: 搜索 FALSE 或 定位
	if !(OperType == "scroll" || strings.TrimSpace(content) == "") {
		isSearch = true
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", PreLineNum, NextLineNum, err
	}

	toal, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, PreLineNum, NextLineNum, err
	}

	//2、定位文件
	location := int64(0)
	if isSearch {
		if "N" == QueryType {
			if NextLineNum == toal {
				err = errors.New("1001:" + "已经搜索到文件结尾")
				return "已经搜索到文件底部", PreLineNum, NextLineNum, err
			}
			location = NextLineNum //向下
		} else if "P" == QueryType {
			if PreLineNum > queryOffSet {
				location = PreLineNum - queryOffSet //分段向上
			} else {
				if "button" == OperType && isSearch {
					err = errors.New("1001:" + "已经搜索到文件开头")
					return "The beginning of the file has been searched", 0, NextLineNum, err
				}
				return nil, 0, NextLineNum, nil
			}
		}
	} else {
		if "N" == QueryType {
			if NextLineNum == toal {
				common.Logger.Info("The end of the file has been searched")
				return nil, PreLineNum, NextLineNum, errors.New("1001:" + "已经搜索到文件结尾")
			}
			location = NextLineNum //向下
		} else if "P" == QueryType {
			if PreLineNum > queryOffSet {
				location = PreLineNum - queryOffSet //分段向上
			} else {
				common.Logger.Info("The end of the file has been searched")
				return nil, 0, PreLineNum, errors.New("1001:" + "已经搜索到文件开头")
			}
		} else if "T" == QueryType {
			if toal-queryOffSet > 0 {
				location = toal - queryOffSet //底部向上
			} else {
				location = 0
			}
			nextOffset = toal
		} else if "H" == QueryType {
			location = 0
		}
	}
	preOffset = location
	offset, err := file.Seek(location, io.SeekStart)
	if err != nil {
		return nil, PreLineNum, NextLineNum, err
	}
	common.Logger.Info("start location ============%d =========", offset)
	defer file.Close()

	//3、搜索函数 /定位函数

	//preOff,nextOff,err := searchFileContent(preQueue,isSearch,content, file,location, toal,QueryType,preOffSet,nextOffset)
	result := make(map[string]interface{})

	searchFileContent2(isSearch, content, file, location, toal, QueryType, result)
	if result["err"] != nil {
		err = result["err"].(error)
		if err != nil {
			common.Logger.Error("search is error :", err)
			return nil, PreLineNum, NextLineNum, err
		}
	}
	preQueue := result["retQueue"].(*queue.Queue)
	preOff := result["preOffSet"].(int64)
	nextOff := result["nextOffset"].(int64)
	findCount := result["findCount"].(int)
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

	defer file.Close()

	/*if findCount < position {
		return "not found 【" + content + "】",returnNum, nil
	}*/

	var buffer = bytes.Buffer{}
	if isSearch && findCount == 0 {
		err = errors.New("1001:" + "没有匹配的关键词【" + content + "】")
	}
	for i := 0; i < preQueue.Length(); i++ {
		buffer.WriteString(preQueue.Get(i).(string))
	}

	return buffer.String(), preOffset, nextOffset, err
}

func searchFileContent2(isSearch bool, search string, file *os.File, offset int64, toal int64, queryType string, result map[string]interface{}) {

	cmp := []byte(search)

	isContinue := false
	retQueue := queue.New()
	findCount := 0
	curOffset := int64(0)
	input := bufio.NewScanner(file)

	isEnd := false
	preOffSet := int64(0)
	nextOffset := int64(0)

	for input.Scan() {
		info := input.Bytes()
		curOffset += int64(len(info)) + 1 //1 代表换行
		if isSearch {

			if bytes.Contains(info, cmp) {
				findCount++
			}
			if findCount == 0 {
				continue
			}
			if retQueue.Length() > 64 {
				retQueue.Remove()
			}
		}
		stext := string(info)
		retQueue.Add(stext + "\n")

		//common.Logger.Info("query count ===%d===%d ===%d == %d ", len(sbyte), curOffset, offset, queryOffSet)

		if curOffset >= queryOffSet2 {

			if findCount > 0 || !isSearch {
				if "P" == queryType {
					preOffSet = offset
					nextOffset = int64(curOffset) + offset
				} else {
					preOffSet = offset
					nextOffset = int64(curOffset) + offset
				}
				isEnd = true
				break
			} else {
				repeat := int64(0)
				if "P" == queryType {
					if offset == 0 {
						isEnd = true
						preOffSet = 0
						break
					}
					if offset-queryOffSet > 0 {
						repeat = offset - queryOffSet
					}
					nextOffset = int64(curOffset) + offset

				} else if "N" == queryType {
					nextOffset = int64(curOffset) + offset
					repeat = nextOffset
				}
				offset, _ = file.Seek(repeat, io.SeekStart)
				isContinue = true
				break
			}
		}
	}

	common.Logger.Info("===================offset===========%d =========", offset)
	if isContinue {
		searchFileContent2(isSearch, search, file, offset, toal, queryType, result)
	} else {
		curoffset, _ := file.Seek(0, io.SeekCurrent)
		common.Logger.Info("===================findCount===========%d ====curoffset %d==", findCount, curoffset)
		if !isEnd {
			preOffSet = offset
			nextOffset = toal
		}
		result["preOffSet"] = preOffSet
		result["nextOffset"] = nextOffset
		result["retQueue"] = retQueue
		result["findCount"] = findCount

	}

}
