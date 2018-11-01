package services

import (
	"bufio"
	"bytes"
	"github.com/eapache/queue"
	"io"
	"logManager/src/common"
	"logManager/src/models"
	"os"
	"strings"
)

const (
	queryOffSet = 10240
)

/**
 * 1\出现次数
 */
func MonitorFileServiceQueryContent(param models.RequestFileParam) (data interface{}, preOff int64, retOffset int64, err error) {

	remoteAddr := strings.TrimSpace(param.RemoteAddr)
	filePath := strings.TrimSpace(param.FilePath)
	content := strings.TrimSpace(param.Content)

	LineNum := param.LineNum
	PreLineNum := param.PreLineNum

	QueryType := param.QueryType
	OperType := param.OperType

	common.Logger.Info(remoteAddr)

	if "N" == QueryType {
		retOffset = LineNum
	} else if "P" == QueryType {
		retOffset = PreLineNum
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, retOffset, err
	}

	toal, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, 0, retOffset, err
	}

	offset := int64(0)
	if "N" == QueryType {

		offset = LineNum
		offset, err = file.Seek(LineNum, io.SeekStart)
		if err != nil {
			return nil, 0, retOffset, err
		}
		common.Logger.Info("query next offset is %d\n", offset)
	} else if "P" == QueryType {

		if PreLineNum > queryOffSet {
			offset = PreLineNum - queryOffSet
		}
		if PreLineNum == 0 {
			common.Logger.Info("query pre is top offset is %d\n", offset)
			return nil, 0, 0, nil
		}

		offset, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, 0, retOffset, err
		}

		common.Logger.Info("query pre offset is %d\n", offset)

	} else if "T" == QueryType {

		if toal >= queryOffSet {
			offset = toal - queryOffSet
		}

		offset, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, 0, retOffset, err
		}

		common.Logger.Info("query TAIL offset is %d\n", offset)

	}

	defer file.Close()

	preQueue, retOffset := searchFileContent(content, file, offset, QueryType, OperType, retOffset, toal)

	/*if findCount < position {
		return "not found 【" + content + "】",returnNum, nil
	}*/

	var buffer = bytes.Buffer{}
	for i := 0; i < preQueue.Length(); i++ {
		buffer.WriteString(preQueue.Get(i).(string))
	}

	return buffer.String(), offset, retOffset, nil
}

func searchFileContent(content string, file *os.File, offset int64, QueryType string, OperType string, retOffset int64, toal int64) (*queue.Queue, int64) {
	isEmpty := true
	if strings.TrimSpace(content) != "" {
		isEmpty = false
	}

	lower := strings.ToLower(content)
	isEnd := false
	isContinue := false
	preQueue := queue.New()
	findCount := 0
	curOffset := 0
	buf := bufio.NewReader(file)
	for {
		stext, err := buf.ReadString('\n')
		if err == io.EOF {
			common.Logger.Info("read the end of file ")
			break
		}
		if err != nil && err != io.EOF {
			common.Logger.Error("read err ", err)
			break
			//return nil,0, retOffset, err
		}
		if !isEmpty {
			text := strings.ToLower(stext)
			if strings.Index(text, lower) >= 0 {
				findCount++
			}
			if findCount == 0 && preQueue.Length() > 10 {
				preQueue.Remove()
			}
		}

		preQueue.Add(stext)

		sbyte := []byte(stext)
		curOffset += len(sbyte) //1 代表换行
		common.Logger.Info("query count ===%d===%d ===%d == %d ", len(sbyte), curOffset, offset, queryOffSet)
		if curOffset >= queryOffSet {
			if isEmpty || OperType == "scroll" { //查询内容为空

				if "P" == QueryType {
					retOffset = offset
				} else {
					retOffset = int64(curOffset) + offset
				}
				isEnd = true
				break

			} else { //查询内容非空

				if findCount > 0 {
					if "P" == QueryType {
						retOffset = offset
					} else {
						retOffset = int64(curOffset) + offset
					}
					isEnd = true
					break
				} else if "P" == QueryType {
					if offset == 0 {
						isEnd = true
						retOffset = 0
						break
					}
					if offset-queryOffSet > 0 {
						retOffset = offset - queryOffSet
					} else {
						retOffset = 0
					}

					offset, err = file.Seek(retOffset, io.SeekStart)
					if err != nil {
						common.Logger.Error("reset err ", err)
						return preQueue, offset
					}
					isContinue = true
					break
				}

			}
		}
	}

	if isContinue {
		searchFileContent(content, file, offset, QueryType, OperType, retOffset, toal)
	} else {
		if !isEnd {
			retOffset = toal
		}
		return preQueue, retOffset
	}
	return nil, 0
}
