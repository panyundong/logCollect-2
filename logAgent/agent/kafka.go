package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type KafkaMessage struct {
	line  string
	topic string
}

type KafkaSend struct {
	client   sarama.SyncProducer
	lineChan chan *KafkaMessage
}

var kafkaSend = &KafkaSend{}

func InitKafka(address []string, threadNum int) (err error) {
	kafkaSend = &KafkaSend{
		lineChan: make(chan *KafkaMessage, 1000),
	}

	kafkaServerConf := sarama.NewConfig()
	kafkaServerConf.Producer.RequiredAcks = sarama.WaitForAll          //ack all
	kafkaServerConf.Producer.Partitioner = sarama.NewRandomPartitioner //random partition
	kafkaServerConf.Producer.Return.Successes = true

	kafkaSend.client, err = sarama.NewSyncProducer(address, kafkaServerConf)
	if err != nil {
		fmt.Println("producer create failed,", err)
		return
	}

	//创建发送的协程
	for i := 0; i < threadNum; i++ {
		logs.Info("start kafka send")
		waitGroup.Add(1)
		go kafkaSend.sendKafka()
	}

	return
}

func (k *KafkaSend) sendKafka() {
	defer waitGroup.Done()
	for value := range k.lineChan {
		message := &sarama.ProducerMessage{
			Topic: value.topic,
			Value: sarama.StringEncoder(value.line),
		}

		partition, offset, err := k.client.SendMessage(message)
		if err != nil {
			logs.Error("send massage to kafka error: %v", err)
			return
		}
		logs.Debug("send message to kafka success. line =%s topic=%s  partition=%s offset=%s", value.line, value.topic, partition, offset)

	}

}

func (k *KafkaSend) addKafkaMessage(line string, topic string) {
	message := &KafkaMessage{
		line:  line,
		topic: topic,
	}
	k.lineChan <- message
	return
}
