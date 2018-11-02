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
	queryOffSet = int64(4096)
)

/**
 * 1\出现次数
 */
func MonitorFileServiceQueryContent(param models.RequestFileParam) (data interface{}, preOffset int64, nextOffset int64, err error) {

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
				return "The end of the file has been searched", PreLineNum, NextLineNum, errors.New("1001")
			}
			location = NextLineNum //向下
		} else if "P" == QueryType {
			if PreLineNum > queryOffSet {
				location = PreLineNum - queryOffSet //分段向上
			} else {
				if "button" == OperType && isSearch {
					return "The beginning of the file has been searched", 0, NextLineNum, errors.New("1001")
				}
				return nil, 0, NextLineNum, nil
			}
		}
	} else {
		if "N" == QueryType {
			if NextLineNum == toal {
				common.Logger.Info("The end of the file has been searched")
				return nil, PreLineNum, NextLineNum, errors.New("1001")
			}
			location = NextLineNum //向下
		} else if "P" == QueryType {
			if PreLineNum > queryOffSet {
				location = PreLineNum - queryOffSet //分段向上
			} else {
				common.Logger.Info("The end of the file has been searched")
				return nil, 0, NextLineNum, errors.New("1001")
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

	searchFileContent(isSearch, content, file, location, toal, QueryType, result)
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
		buffer.WriteString("not found 【" + content + "】")
	} else {
		for i := 0; i < preQueue.Length(); i++ {
			buffer.WriteString(preQueue.Get(i).(string))
		}
	}

	return buffer.String(), preOffset, nextOffset, nil
}

func searchFileContent(isSearch bool, search string, file *os.File, offset int64, toal int64, queryType string, result map[string]interface{}) {

	lower := strings.ToLower(search)
	isContinue := false
	retQueue := queue.New()
	findCount := 0
	curOffset := int64(0)
	buf := bufio.NewReader(file)

	isEnd := false
	preOffSet := int64(0)
	nextOffset := int64(0)

	for {
		line, isPrefix, err := buf.ReadLine()
		if err == io.EOF {
			common.Logger.Info("isPrefix:", isPrefix)
			nextOffset = int64(curOffset) + offset
			common.Logger.Info("read the end of file ")
			result["err"] = nil
			break
		}
		if err != nil && err != io.EOF {
			common.Logger.Error("read err ", err)
			result["err"] = err
			break
		}
		stext := string(line)
		if isSearch {
			text := strings.ToLower(stext)
			if strings.Index(text, lower) >= 0 {
				findCount++
			}
			if findCount == 0 && retQueue.Length() > 50 {
				retQueue.Remove()
			}
		}

		retQueue.Add(stext)

		curOffset += int64(len(line)) //1 代表换行
		//common.Logger.Info("query count ===%d===%d ===%d == %d ", len(sbyte), curOffset, offset, queryOffSet)

		if curOffset >= queryOffSet {

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
				offset, err = file.Seek(repeat, io.SeekStart)
				if err != nil {
					common.Logger.Error("reset err ", err)
					result["err"] = err
					break
				}
				isContinue = true
				break
			}
		}
	}

	common.Logger.Info("===================offset===========%d ======", offset)
	if isContinue {
		searchFileContent(isSearch, search, file, offset, toal, queryType, result)
	} else {
		if !isEnd {
			preOffSet = offset
			nextOffset = toal
		}
		result["preOffSet"] = preOffSet
		result["nextOffset"] = nextOffset
		result["retQueue"] = retQueue
		result["findCount"] = findCount

		//return preOffSet,nextOffset,nil
	}

	//return 0,nextOffset,nil
}
