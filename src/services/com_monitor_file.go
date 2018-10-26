package services

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/eapache/queue"
	"logManager/src/common"
	"logManager/src/models"
	"os"
	"strings"
)

const (
	queryCount = 100
)

/**
 * 1\出现次数
 */
func MonitorFileServiceQueryContent(param models.RequestFileParam) (data interface{}, returnNum int, err error) {

	remoteAddr := strings.TrimSpace(param.RemoteAddr)
	filePath := strings.TrimSpace(param.FilePath)
	content := strings.TrimSpace(param.Content)
	LineNum := param.LineNum
	PreLineNum := param.PreLineNum
	position := param.Position
	fmt.Println(position)
	QueryType := param.QueryType

	common.Logger.Info(remoteAddr)

	if "N" == QueryType {
		returnNum = LineNum
	} else if "P" == QueryType {
		returnNum = PreLineNum
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", returnNum, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if err != nil {
		return nil, returnNum, err
	}
	preQueue := queue.New()
	//nextQueue := queue.New()
	findCount := 0

	lcontent := strings.ToLower(content)

	lineCount := 1
	minNum := 0
	maxNum := 0
	if "N" == QueryType {
		maxNum = LineNum + queryCount
		minNum = LineNum
	} else if "P" == QueryType {
		maxNum = PreLineNum
		minNum = PreLineNum - queryCount
	}

	if err := scanner.Err(); err != nil {
		return "", returnNum, err
	}

	for scanner.Scan() {
		if (lineCount >= minNum && lineCount < maxNum) || QueryType == "T" { //T 文件尾部
			stext := scanner.Text()
			text := strings.ToLower(stext)
			if strings.Index(text, lcontent) >= 0 {
				findCount++
			}
			//if findCount == position-1 {
			preQueue.Add(stext)
			//}

			if preQueue.Length() > queryCount {
				preQueue.Remove()
			}

			/*if findCount >= position && nextQueue.Length() < 20 {
				if findCount == position {
					nextQueue.Add("<span style='color:red'>" + stext + "</span>")
				} else {
					nextQueue.Add(stext)
				}
			}
			if nextQueue.Length() >= 20 {
				break
			}*/
		}

		lineCount++
		if lineCount >= maxNum && QueryType != "T" {
			break
		}
	}

	if "N" == QueryType {
		returnNum = LineNum + (lineCount - 1)
	} else if "P" == QueryType {
		returnNum = PreLineNum - (lineCount - 2)
	} else if "T" == QueryType {
		returnNum = lineCount - 1
	}

	/*if findCount < position {
		return "not found 【" + content + "】",returnNum, nil
	}*/

	var buffer = bytes.Buffer{}
	for i := 0; i < preQueue.Length(); i++ {
		buffer.WriteString(preQueue.Get(i).(string) + "\n")
	}

	/*for i := 0; i < nextQueue.Length(); i++ {
		buffer.WriteString(nextQueue.Get(i).(string) + "\n")
	}*/

	return buffer.String(), returnNum, nil
}
