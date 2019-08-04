package test

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestKafka(t *testing.T) {
	kafkaServerConf := sarama.NewConfig()
	kafkaServerConf.Producer.RequiredAcks = sarama.WaitForAll
	kafkaServerConf.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaServerConf.Producer.Return.Successes = true

	kafkaServerClient, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, kafkaServerConf)
	if err != nil {
		fmt.Println("producer create failed,", err)
		return
	}

	logs.Info(kafkaServerClient)
}
