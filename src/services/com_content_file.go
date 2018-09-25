package services

import (
	"bufio"
	"bytes"
	"github.com/eapache/queue"
	"logManager/src/common"
	"logManager/src/models"
	"os"
	"strings"
)

/**
 * 1\出现次数
 */
func LogFileServiceQueryContent(param models.RequestFileParam) (data interface{}, err error) {

	remoteAddr := strings.TrimSpace(param.RemoteAddr)
	filePath := strings.TrimSpace(param.FilePath)
	content := strings.TrimSpace(param.Content)
	position := param.Position

	common.Logger.Info(remoteAddr)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if err != nil {
		return nil, err
	}
	preQueue := queue.New()
	nextQueue := queue.New()
	findCount := 0

	lcontent := strings.ToLower(content)
	for scanner.Scan() {
		stext := scanner.Text()
		text := strings.ToLower(stext)
		if strings.Index(text, lcontent) >= 0 {
			findCount++
		}
		if findCount == position-1 {
			preQueue.Add(stext)
		}

		if preQueue.Length() > 20 {
			preQueue.Remove()
		}

		if findCount >= position && nextQueue.Length() < 20 {
			if findCount == position {
				nextQueue.Add("<span style='color:red'>" + stext + "</span>")
			} else {
				nextQueue.Add(stext)
			}
		}
		if nextQueue.Length() >= 20 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if findCount < position {
		return "not found 【" + content + "】", nil
	}
	var buffer = bytes.Buffer{}
	for i := 0; i < preQueue.Length(); i++ {
		buffer.WriteString(preQueue.Get(i).(string) + "\n")
	}

	for i := 0; i < nextQueue.Length(); i++ {
		buffer.WriteString(nextQueue.Get(i).(string) + "\n")
	}

	return buffer.String(), nil
}

//远程调用接口
/*func LogFileServiceQueryContent_remote(remoteAddr string, filePath string,content string) {

	chanName := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	msgKey := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)

	models.Clients[socketConn] = models.NewBroadCast()
	models.BroadCastMap[msgKey] = models.Clients[socketConn] //key->broadCast
	go models.HandlerMessage(msgKey, *socketConn)

	vModel := models.ConfigRemote{}
	vModel.RemoteAddr = remoteAddr + "/open/logFile/startTail"
	vModel.Method = "POST"
	vModel.Header = "{}"
	vModel.Param = "{}"
	vModel.Body = "{\"filePath\":\"" + filePath + "\",\"chanName\":\"" + chanName + "\",\"msgKey\":\"" + msgKey + "\"}"

	remoteServiceStart(vModel)

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object
		err := socketConn.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson ... : %v ", err)
			delete(models.Clients, socketConn)
			defer socketConn.Close()
			remoteServiceStop(chanName, msgKey)

			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}

}*/
