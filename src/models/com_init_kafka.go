package models

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"logManager/src/common"
	"strings"
	"sync"
)

var (
	produce       sarama.SyncProducer
	consumer      sarama.Consumer
	partitionList []int32
	wg            sync.WaitGroup
)

func InitKafka() (err error) {

	enableKafka := beego.AppConfig.String("tailf.kafka.enable")
	if "n" != enableKafka {
		remoteTail := beego.AppConfig.String("tailf.kafka.type")
		if "remote" != remoteTail {
			err = initProduce()
			if err != nil {
				return err
			}
		}

		err = initConsumer()

		if err != nil {
			return err
		}
	}

	return
}

// 初始化KAFKA生产者
func initProduce() (err error) {
	//初始化KAFKA配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5000

	//创建生产者
	kafkaServer := beego.AppConfig.String("kafka.producer.servers")
	servers := strings.Split(strings.TrimSpace(kafkaServer), ",")
	produce, err = sarama.NewSyncProducer(servers, config)

	if err != nil {
		common.Logger.Error("sarama.NewSyncProducer failed ", err)
		return
	}
	return err
}

//初始化消费者
func initConsumer() (err error) {

	//创建消费者
	kafkaServer := beego.AppConfig.String("kafka.consumer.servers")
	servers := strings.Split(strings.TrimSpace(kafkaServer), ",")
	//初始化KAFKA配置
	config := sarama.NewConfig()
	//config.Consumer.MaxWaitTime = 5000
	consumer, err = sarama.NewConsumer(servers, config)
	if err != nil {
		common.Logger.Error("sarama.NewConsumer failed ", err)
		return
	}

	//设置分区
	partitionList, err = consumer.Partitions(common.TopicLog)
	if err != nil {
		common.Logger.Error("Failed to get the list of partitions :", err)
		return
	}
	go SendKafkaMsg2Chan()

	return err
}

func SendToKafka(msgKey, data string) (err error) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = common.TopicLog
	msg.Value = sarama.StringEncoder(data)
	msg.Key = sarama.StringEncoder(msgKey)

	pid, offset, err := produce.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed, err:%v pid:%v offset:%v data:%v topic:%v", err, pid, offset, data)
		return
	}

	return
}

//客户端接收消息发送至chain
func SendKafkaMsg2Chan() {

	//循环分区
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(common.TopicLog, int32(partition), sarama.OffsetNewest)
		if err != nil {
			common.Logger.Error("Failed to start consumer for partition %d : %s \n", partition, err)
			continue
		}

		defer pc.AsyncClose()

		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				msgKey := string(msg.Key)
				//msgText := Message{string(msg.Value)}
				broadCast := BroadCastMap[msgKey] //获取缓存消息通道
				if broadCast.msgChan != nil {
					broadCast.msgChan <- string(msg.Value)
				}
			}

		}(pc)
	}
	wg.Wait()
	common.Logger.Info("Done consuming topic " + common.TopicLog)
	//	consumer.Close()
}
