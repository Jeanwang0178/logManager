package utils

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"logManager/models"
	"sync"
)

var (
	produce       sarama.SyncProducer
	consumer      sarama.Consumer
	partitionList []int32
	wg            sync.WaitGroup
)

func InitKafka() (err error) {

	err = initProduce()
	if err != nil {
		return err
	}

	err = initConsumer()

	if err != nil {
		return err
	}

	return
}

// 初始化KAFKA生产者
func initProduce() (err error) {
	//初始化KAFKA配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//创建生产者
	produce, err = sarama.NewSyncProducer([]string{"192.168.3.186:9092"}, config)

	if err != nil {
		Logger.Error("sarama.NewSyncProducer failed ", err)
		return
	}
	return err
}

//初始化消费者
func initConsumer() (err error) {

	//创建消费者
	consumer, err = sarama.NewConsumer([]string{"192.168.3.186:9092"}, nil)
	if err != nil {
		Logger.Error("sarama.NewConsumer failed ", err)
		return
	}

	//设置分区
	partitionList, err = consumer.Partitions(TopicLog)
	if err != nil {
		Logger.Error("Failed to get the list of partitions :", err)
		return
	}
	//发送消息至chan
	go sendKafkaMsg2Chan()
	return err
}

func SendToKafka(data, topic string) (err error) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	pid, offset, err := produce.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed, err:%v data:%v topic:%v", err, data, topic)
		return
	}

	logs.Debug("send succ, pid:%v offset:%v, topic:%v\n", pid, offset, topic)
	return
}

//客户端接收消息发送至chain
func sendKafkaMsg2Chan() {

	//循环分区
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(TopicLog, int32(partition), sarama.OffsetNewest)
		if err != nil {
			Logger.Error("Failed to start consumer for partition %d : %s \n", partition, err)
			return
		}

		defer pc.AsyncClose()

		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				Logger.Info("partition :%d ,Offset:%d ,key:%s,Value :%s ", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				msg := models.Message{string(msg.Value)}
				Broadcast <- msg
			}
		}(pc)
	}
	wg.Wait()
	Logger.Info("Done consuming topic " + TopicLog)
	consumer.Close()
}
